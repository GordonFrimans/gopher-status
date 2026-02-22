package storage

import (
	"errors"
	"sync"

	"golang.org/x/crypto/bcrypt"
)

// Login = role
// Любой логин кроме admin это обычный пользователь
type User struct {
	Login string
	Hash  []byte
}

type InMemoryUser struct {
	// RWMutex эффективнее, так как позволяет читать (Get)
	// одновременно из нескольких горутин, блокируя только при записи
	mu    sync.RWMutex
	Users map[string]*User
}

// При инициализации сразу загоняем данные админа
func NewInMemoryUser() *InMemoryUser {
	inmemoryUer := &InMemoryUser{
		Users: make(map[string]*User),
	}
	// Игнорируем ошибку здесь (т.к. мы уверены в хардкод-данных),
	// но можно обернуть в панику, если админ не создался
	_ = inmemoryUer.Create("admin", "Votomaj200700")
	return inmemoryUer
}

func (u *InMemoryUser) Login(login, password string) (bool, error) {
	u.mu.RLock()
	defer u.mu.RUnlock()

	if _, ok := u.Users[login]; !ok {
		return false, errors.New("Неверный логин или пароль")
	}
	return u.Users[login].CheckPassword(password), nil
}

func (u *InMemoryUser) Create(login, password string) error {
	u.mu.Lock()
	defer u.mu.Unlock()

	if _, ok := u.Users[login]; ok {
		return errors.New("Login Busy")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.Users[login] = &User{
		Login: login,
		Hash:  hash,
	}
	return nil
}

// GetUser
func (u *InMemoryUser) Get(login string) (User, error) {
	u.mu.RLock() // Используем RLock (Read Lock) для параллельного чтения
	defer u.mu.RUnlock()

	userPtr, ok := u.Users[login]
	if !ok {
		return User{}, errors.New("user not found")
	}

	// Возвращаем значение по указателю (копию),
	// чтобы извне не могли случайно изменить хэш в мапе
	return *userPtr, nil
}

// DeleteUser
func (u *InMemoryUser) Delete(login string) error {
	u.mu.Lock()
	defer u.mu.Unlock()

	// Защита от "выстрела в ногу": запрещаем удалять главного админа
	if login == "admin" {
		return errors.New("cannot delete root admin")
	}

	if _, ok := u.Users[login]; !ok {
		return errors.New("user not found")
	}

	delete(u.Users, login)
	return nil
}

// Полезный хелпер: проверяет, совпадает ли пароль с хешем пользователя.
// Очень пригодится тебе в gRPC-хендлере Login() !
func (user *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword(user.Hash, []byte(password))
	return err == nil
}
