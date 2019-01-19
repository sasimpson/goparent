package datastore

import (
	"context"
	"crypto/md5"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/sasimpson/goparent"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/mail"
)

//UserService -
type UserService struct {
	Env *goparent.Env
}

//UserClaims -
type UserClaims struct {
	ID       string
	Name     string
	Email    string
	Username string
	Password string
	jwt.StandardClaims
}

//UserKind - constant string for all user entities in datastore
const UserKind = "User"

//ResetKind - constant for all user password reset entities
const ResetKind = "PasswordReset"

var (
	//ErrNoUserFound is when there is no user returned by the datastore
	ErrNoUserFound = errors.New("no result for that id")
	//ErrInvalidLogin is when the password/user combo do not match
	ErrInvalidLogin = errors.New("no result for that username password combo")
	//ErrInvalidEmail is when a user submits an invalid email for password reset
	ErrInvalidEmail = errors.New("no result for that email")
	//ErrInvalidResetCode is when a user submits a code for resetting password that is invalid
	ErrInvalidResetCode = errors.New("invalid code for reset")
)

//User - get a user by the key/id
func (s *UserService) User(ctx context.Context, key string) (*goparent.User, error) {
	var user goparent.User
	userKey := datastore.NewKey(ctx, UserKind, key, 0, nil)
	err := datastore.Get(ctx, userKey, &user)
	if err != nil {
		if err == datastore.ErrNoSuchEntity {
			return nil, NewError("datastore.UserService.User", ErrNoUserFound)
		}
		return nil, err
	}
	return &user, nil
}

//UserByLogin - get a user by their login, email and password
func (s *UserService) UserByLogin(ctx context.Context, email string, password string) (*goparent.User, error) {
	var user goparent.User
	userKey := datastore.NewKey(ctx, UserKind, md5Email(email), 0, nil)
	err := datastore.Get(ctx, userKey, &user)
	if err != nil {
		if err == datastore.ErrNoSuchEntity {
			return nil, ErrInvalidLogin
		}
		return nil, NewError("datastore.UserService.User", err)
	}

	//verify the email and password match, then return it
	if user.Password == password && user.Email == email {
		return &user, nil
	}
	return nil, ErrInvalidLogin
}

//Save - save a user's current values
func (s *UserService) Save(ctx context.Context, user *goparent.User) error {
	userKey := datastore.NewKey(ctx, UserKind, md5Email(user.Email), 0, nil)
	//we use id as a way to lookup stuff, but don't actually _use_ id since we are using key...

	var family *goparent.Family
	fs := &FamilyService{Env: s.Env}
	//user doesn't have a current family
	if user.CurrentFamily == "" && user.ID != "" {
		//get the family for which the user is the admin
		family, err := fs.GetAdminFamily(ctx, user)
		if err == ErrNoFamilyFound {
			//didn't find a family, creating one, save it
			family = &goparent.Family{Admin: userKey.StringID(), Members: []string{userKey.StringID()}}
			err = fs.Save(ctx, family)
			if err != nil {
				return NewError("datastore.UserService.Save.1a", err)
			}
		}
		if err != nil {
			return NewError("datastore.UserService.Save.1b", err)
		}
		user.CurrentFamily = family.ID
	}

	//if the user has no current family and no id, then we need to create a family and save it
	if user.CurrentFamily == "" && user.ID == "" {
		family = &goparent.Family{Admin: userKey.StringID(), Members: []string{userKey.StringID()}}
		err := fs.Save(ctx, family)
		if err != nil {
			return NewError("datastore.UserService.Save.2", err)
		}
		user.CurrentFamily = family.ID
	}

	user.ID = userKey.StringID()
	//this will save if there is or isn't a record for this user.
	_, err := datastore.Put(ctx, userKey, user)
	if err != nil {
		return NewError("datastore.UserService.Save.3", err)
	}
	return nil
}

//GetToken - gets the user token
func (s *UserService) GetToken(user *goparent.User, duration time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["Name"] = user.Name
	claims["ID"] = user.ID
	claims["Email"] = user.Email
	claims["Username"] = user.Username
	claims["exp"] = time.Now().Add(duration).Unix()
	tokenString, err := token.SignedString(s.Env.Auth.SigningKey)
	if err != nil {
		return "", NewError("datastore.UserService.GetToken", err)
	}
	return tokenString, nil
}

//ValidateToken - validate token against signing method and populate user.
func (s *UserService) ValidateToken(ctx context.Context, tokenString string) (*goparent.User, bool, error) {
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return s.Env.Auth.SigningKey, nil
	})
	if err != nil {
		return nil, false, NewError("datastore.UserService.ValidateToken", err)
	}

	if claims, ok := token.Claims.(*UserClaims); ok && token.Valid {
		user, err := s.User(ctx, claims.ID)
		if err != nil {
			return nil, false, NewError("datastore.UserService.ValidateToken", err)
		}
		return user, true, nil
	}
	return nil, false, errors.New("invalid token")
}

//GetFamily -
func (s *UserService) GetFamily(ctx context.Context, user *goparent.User) (*goparent.Family, error) {
	if user.CurrentFamily == "" {
		return nil, errors.New("user has no current family")
	}

	fs := FamilyService{Env: s.Env}
	family, err := fs.Family(ctx, user.CurrentFamily)
	if err != nil {
		return nil, NewError("datastore.UserService.GetFamily", err)
	}
	return family, nil
}

//GetAllFamily -
func (s *UserService) GetAllFamily(ctx context.Context, user *goparent.User) ([]*goparent.Family, error) {
	var families []*goparent.Family
	q := datastore.NewQuery(FamilyKind).Filter("Members =", user.ID)
	t := q.Run(ctx)
	for {
		var family goparent.Family
		_, err := t.Next(&family)
		if err == datastore.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		families = append(families, &family)
	}
	return families, nil
}

//RequestResetPassword will setup a password reset token for the email submitted.  that token will then be used
// to actually change the password
func (s *UserService) RequestResetPassword(ctx context.Context, email string, ip string) error {
	q := datastore.NewQuery(UserKind).Filter("Email =", email)
	t := q.Run(ctx)
	var user goparent.User
	_, err := t.Next(&user)
	if err != nil {
		return ErrInvalidEmail
	}

	// log.Println("password reset requested for user", user.Email)

	// userKey := datastore.NewKey(ctx, UserKind, user.ID, 0, nil)
	resetKey := datastore.NewIncompleteKey(ctx, ResetKind, nil)

	resetRequest := &goparent.UserReset{
		Timestamp:   time.Now(),
		RequestAddr: ip,
		Email:       email,
	}

	key, err := datastore.Put(ctx, resetKey, resetRequest)
	if err != nil {
		return err
	}
	//this key needs to be emailed to the user.  should eventually make this a jwt reset.
	// log.Printf("key: %#v", key)
	code := encodeInt(key.IntID())
	// log.Println("password reset key is", key.IntID(), code)

	resetMessage := mail.Message{
		Sender:  "noreply@goparent-181120.appspotmail.com",
		To:      []string{user.Email},
		Subject: "GoParent password reset",
		Body:    fmt.Sprintf("your password reset code is: %s", code),
	}

	err = mail.Send(ctx, &resetMessage)
	if err != nil {
		log.Printf("error sending mail: %#v", err)
	}

	return nil
}

//ResetPassword will reset the password for the user assuming they meet the requirements
func (s *UserService) ResetPassword(ctx context.Context, code string, password string) error {
	//get code and verify it exists in the datastore
	resetID := decodeBase64(code)
	log.Printf("code: %s", code)
	log.Printf("id: %d", resetID)
	resetKey := datastore.NewKey(ctx, ResetKind, "", resetID, nil)
	var resetRequest goparent.UserReset
	err := datastore.Get(ctx, resetKey, &resetRequest)
	if err == datastore.ErrNoSuchEntity {
		return ErrInvalidResetCode
	}
	if err != nil {
		return NewError("datastore.ResetPassword a", err)
	}

	//now lookup the user and reset the password to the new one.
	q := datastore.NewQuery(UserKind).Filter("Email =", resetRequest.Email)
	t := q.Run(ctx)
	var user goparent.User
	_, err = t.Next(&user)
	if err != nil {
		return NewError("datastore.ResetPassword b", err)
	}

	userKey := datastore.NewKey(ctx, UserKind, md5Email(user.Email), 0, nil)
	user.Password = password
	_, err = datastore.Put(ctx, userKey, &user)
	if err != nil {
		return NewError("datastore.ResetPassword c", err)
	}

	err = datastore.Delete(ctx, resetKey)
	if err != nil {
		return NewError("datastore.ResetPassword d", err)
	}

	return nil
}

//util functions
func md5Email(email string) string {
	h := md5.New()
	h.Write([]byte(email))
	return hex.EncodeToString(h.Sum(nil))
}

func encodeInt(i int64) string {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(i))
	return base64.StdEncoding.EncodeToString(b)
}

func decodeBase64(encoded string) int64 {
	decoded, _ := base64.StdEncoding.DecodeString(encoded)
	return int64(binary.LittleEndian.Uint64(decoded))
}
