package main

import (
	"math"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

type Option func(*GamePerson)

func WithName(name string) func(*GamePerson) {
	return func(person *GamePerson) {
		copy(person.name[:], name)
	}
}

func WithCoordinates(x, y, z int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.x = int32(x)
		person.y = int32(y)
		person.z = int32(z)
	}
}

func WithGold(gold int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.gold = int32(gold)
	}
}

func WithMana(mana int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.mana = uint16(mana)
	}
}

func WithHealth(health int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.health = uint16(health)
	}
}

func WithRespect(respect int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.respect = uint8(respect)
	}
}

func WithStrength(strength int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.strength = uint8(strength)
	}
}

func WithExperience(experience int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.experience = uint8(experience)
	}
}

func WithLevel(level int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.level = uint8(level)
	}
}

func WithHouse() func(*GamePerson) {
	return func(person *GamePerson) {
		person.flags |= flagHasHouse
	}
}

func WithGun() func(*GamePerson) {
	return func(person *GamePerson) {
		person.flags |= flagHasGun
	}
}

func WithFamily() func(*GamePerson) {
	return func(person *GamePerson) {
		person.flags |= flagHasFamily
	}
}

func WithType(personType int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.personType = uint8(personType)
	}
}

const (
	BuilderGamePersonType = iota
	BlacksmithGamePersonType
	WarriorGamePersonType
)

// Битовые маски и сдвиги для упакованных данных
const (
	respectMask     = 0x0F       // 4 bits (0-15)
	strengthMask    = 0x0F       // 4 bits
	experienceMask  = 0x0F       // 4 bits
	levelMask       = 0x0F       // 4 bits
	personTypeMask  = 0x0F       // 4 bits
	flagsMask       = 0x07       // 3 bits
	
	respectShift    = 0
	strengthShift   = 4
	experienceShift = 8
	levelShift      = 12
	personTypeShift = 16
	flagsShift      = 20
)

// Флаги для битовых операций
const (
	flagHasHouse  = 1 << 0 // 0001
	flagHasGun    = 1 << 1 // 0010
	flagHasFamily = 1 << 2 // 0100
)

// GamePerson - оптимизированная структура размером 64 байта
type GamePerson struct {
	// Координаты и золото - int32 (16 bytes)
	x, y, z, gold int32 // 16 bytes
	
	// Мана и здоровье - uint16 (4 bytes)  
	mana, health uint16 // 4 bytes
	
	// Имя пользователя - 40 байт (вместо 42 для экономии места)
	name [40]byte // 40 bytes
	
	// Упаковываем мелкие значения в один uint32 битовыми операциями
	// Respect (4 bits, max 15), Strength (4 bits), Experience (4 bits), 
	// Level (4 bits), PersonType (4 bits), Flags (3 bits) = 23 bits total
	packed uint32 // 4 bytes
	
	// Итого: 16 + 4 + 40 + 4 = 64 bytes точно!
}

func NewGamePerson(options ...Option) GamePerson {
	person := GamePerson{}
	for _, option := range options {
		option(&person)
	}
	return person
}

func (p *GamePerson) Name() string {
	// Находим длину строки (до первого нулевого байта)
	for i, b := range p.name {
		if b == 0 {
			return string(p.name[:i])
		}
	}
	return string(p.name[:])
}

func (p *GamePerson) X() int {
	return int(p.x)
}

func (p *GamePerson) Y() int {
	return int(p.y)
}

func (p *GamePerson) Z() int {
	return int(p.z)
}

func (p *GamePerson) Gold() int {
	return int(p.gold)
}

func (p *GamePerson) Mana() int {
	return int(p.mana)
}

func (p *GamePerson) Health() int {
	return int(p.health)
}

func (p *GamePerson) Respect() int {
	return int(p.respect)
}

func (p *GamePerson) Strength() int {
	return int(p.strength)
}

func (p *GamePerson) Experience() int {
	return int(p.experience)
}

func (p *GamePerson) Level() int {
	return int(p.level)
}

func (p *GamePerson) HasHouse() bool {
	return p.flags&flagHasHouse != 0
}

func (p *GamePerson) HasGun() bool {
	return p.flags&flagHasGun != 0
}

func (p *GamePerson) HasFamilty() bool {
	return p.flags&flagHasFamily != 0
}

func (p *GamePerson) Type() int {
	return int(p.personType)
}

func TestGamePerson(t *testing.T) {
	assert.LessOrEqual(t, unsafe.Sizeof(GamePerson{}), uintptr(64))

	const x, y, z = math.MinInt32, math.MaxInt32, 0
	const name = "aaaaaaaaaaaaa_bbbbbbbbbbbbb_cccccccccccccc"
	const personType = BuilderGamePersonType
	const gold = math.MaxInt32
	const mana = 1000
	const health = 1000
	const respect = 10
	const strength = 10
	const experience = 10
	const level = 10

	options := []Option{
		WithName(name),
		WithCoordinates(x, y, z),
		WithGold(gold),
		WithMana(mana),
		WithHealth(health),
		WithRespect(respect),
		WithStrength(strength),
		WithExperience(experience),
		WithLevel(level),
		WithHouse(),
		WithFamily(),
		WithType(personType),
	}

	person := NewGamePerson(options...)
	assert.Equal(t, name, person.Name())
	assert.Equal(t, x, person.X())
	assert.Equal(t, y, person.Y())
	assert.Equal(t, z, person.Z())
	assert.Equal(t, gold, person.Gold())
	assert.Equal(t, mana, person.Mana())
	assert.Equal(t, health, person.Health())
	assert.Equal(t, respect, person.Respect())
	assert.Equal(t, strength, person.Strength())
	assert.Equal(t, experience, person.Experience())
	assert.Equal(t, level, person.Level())
	assert.True(t, person.HasHouse())
	assert.True(t, person.HasFamilty())
	assert.False(t, person.HasGun())
	assert.Equal(t, personType, person.Type())
}