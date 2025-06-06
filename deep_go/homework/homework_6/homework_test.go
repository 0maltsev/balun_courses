package main

import (
	"math"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

const (
	manaMask        = 0x03FF
	personTypeShift = 10

	healthShift     = 3
	houseBit        = 0x0001
	gunBit          = 0x0002
	familyBit       = 0x0004

	respectShift    = 12
	strengthShift   = 8
	experienceShift = 4
)

type Option func(*GamePerson)

func WithName(name string) Option {
	return func(person *GamePerson) {
		n := copy(person.name[:], name)
		for i := n; i < len(person.name); i++ {
			person.name[i] = 0
		}
	}
}

func WithCoordinates(x, y, z int) Option {
	return func(person *GamePerson) {
		person.x, person.y, person.z = int32(x), int32(y), int32(z)
	}
}

func WithGold(gold int) Option {
	return func(person *GamePerson) {
		person.gold = int32(gold)
	}
}

func WithMana(mana int) Option {
	return func(person *GamePerson) {
		person.setMana(mana)
	}
}

func WithType(personType int) Option {
	return func(person *GamePerson) {
		person.setType(personType)
	}
}

func WithHealth(health int) Option {
	return func(person *GamePerson) {
		person.setHealth(health)
	}
}

func WithHouse() Option {
	return func(person *GamePerson) {
		person.setHouse(true)
	}
}

func WithGun() Option {
	return func(person *GamePerson) {
		person.setGun(true)
	}
}

func WithFamily() Option {
	return func(person *GamePerson) {
		person.setFamily(true)
	}
}

func WithRespect(respect int) Option {
	return func(person *GamePerson) {
		person.setRespect(respect)
	}
}

func WithStrength(strength int) Option {
	return func(person *GamePerson) {
		person.setStrength(strength)
	}
}

func WithExperience(exp int) Option {
	return func(person *GamePerson) {
		person.setExperience(exp)
	}
}

func WithLevel(level int) Option {
	return func(person *GamePerson) {
		person.setLevel(level)
	}
}

const (
	BuilderGamePersonType = iota
	BlacksmithGamePersonType
	WarriorGamePersonType
)

type GamePerson struct {
	x, y, z  int32
	gold     int32
	name     [42]byte

	// packed fields
	levelExpirienceStrengthRespect uint16
	manaType                       int16
	healthHouseGunFamily           int16
}

func NewGamePerson(options ...Option) GamePerson {
	p := GamePerson{}
	for _, opt := range options {
		opt(&p)
	}
	return p
}

func (p *GamePerson) setMana(mana int) {
	p.manaType = (p.manaType &^ int16(manaMask)) | int16(mana&manaMask)
}

func (p *GamePerson) setType(personType int) {
	p.manaType = (p.manaType & int16(manaMask)) | int16(personType<<personTypeShift)
}

func (p *GamePerson) setHealth(health int) {
	p.healthHouseGunFamily = (p.healthHouseGunFamily & 0x0007) | int16(health<<healthShift)
}

func (p *GamePerson) setHouse(on bool) {
	if on {
		p.healthHouseGunFamily |= houseBit
	} else {
		p.healthHouseGunFamily &^= houseBit
	}
}

func (p *GamePerson) setGun(on bool) {
	if on {
		p.healthHouseGunFamily |= gunBit
	} else {
		p.healthHouseGunFamily &^= gunBit
	}
}

func (p *GamePerson) setFamily(on bool) {
	if on {
		p.healthHouseGunFamily |= familyBit
	} else {
		p.healthHouseGunFamily &^= familyBit
	}
}

func (p *GamePerson) setRespect(val int) {
	mask := uint16(0xF << respectShift)
	p.levelExpirienceStrengthRespect = (p.levelExpirienceStrengthRespect &^ mask) | (uint16(val&0xF) << respectShift)
}

func (p *GamePerson) setStrength(val int) {
	mask := uint16(0xF << strengthShift)
	p.levelExpirienceStrengthRespect = (p.levelExpirienceStrengthRespect &^ mask) | (uint16(val&0xF) << strengthShift)
}

func (p *GamePerson) setExperience(val int) {
	mask := uint16(0xF << experienceShift)
	p.levelExpirienceStrengthRespect = (p.levelExpirienceStrengthRespect &^ mask) | (uint16(val&0xF) << experienceShift)
}

func (p *GamePerson) setLevel(val int) {
	mask := uint16(0xF)
	p.levelExpirienceStrengthRespect = (p.levelExpirienceStrengthRespect &^ mask) | uint16(val&0xF)
}

func (p *GamePerson) Name() string {
	length := 0
	for i, b := range p.name {
		if b == 0 {
			break
		}
		length = i + 1
	}
	return unsafe.String(&p.name[0], length)
}

func (p *GamePerson) X() int              { return int(p.x) }
func (p *GamePerson) Y() int              { return int(p.y) }
func (p *GamePerson) Z() int              { return int(p.z) }
func (p *GamePerson) Gold() int           { return int(p.gold) }
func (p *GamePerson) Mana() int           { return int(p.manaType & manaMask) }
func (p *GamePerson) Type() int           { return int(p.manaType >> personTypeShift) }
func (p *GamePerson) Health() int         { return int(p.healthHouseGunFamily>>healthShift) & 0x07FF }
func (p *GamePerson) HasHouse() bool      { return p.healthHouseGunFamily&houseBit != 0 }
func (p *GamePerson) HasGun() bool        { return p.healthHouseGunFamily&gunBit != 0 }
func (p *GamePerson) HasFamily() bool     { return p.healthHouseGunFamily&familyBit != 0 }
func (p *GamePerson) Respect() int        { return int(p.levelExpirienceStrengthRespect>>respectShift) & 0xF }
func (p *GamePerson) Strength() int       { return int(p.levelExpirienceStrengthRespect>>strengthShift) & 0xF }
func (p *GamePerson) Experience() int     { return int(p.levelExpirienceStrengthRespect>>experienceShift) & 0xF }
func (p *GamePerson) Level() int          { return int(p.levelExpirienceStrengthRespect & 0xF) }


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
	assert.True(t, person.HasFamily())
	assert.False(t, person.HasGun())
	assert.Equal(t, personType, person.Type())
}