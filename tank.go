// tank.go
package main

import (
	"math"
	"math/rand"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

/*

It is not easy, to create this file to follow only one function set.
Let's say: tank.moveForward(angle float64) -> this seems to be taking the tank pointer, and change it's angle.
But in my case, it will take a value receiver, change it's angle and return that. ((a tank) func moveForwardangle float64) tank)
this is because I need to check whether the tank is overlapping with other game objects or not.
This checking(or rather collision detection) has been done in the main.go file, because that is the central place, where I
have acces to all game objects, so I have to essentially return a new tank, and not mutate the receiver.

In some cases the function signature makes it very clear that it will mutate, for those cases I have used pointer receivers
like bullet.Update()

func (tank *EnemyTank) WillUpdate() bool -> this is exceptional, it looks like it's not mutating but it needs to do that
see it's usage in main.go, there is no better or elegant way of doing it in any other way(maybe)

callbacks need to be pointers, value receivers will 'bind' the function with the receiver. Changing the receiver is not possible.(see player update in main.go)


*/

type Bullet struct {
	bulletTexture *sdl.Texture
	velocity      float32
	boundingBox   sdl.FRect
	rotationAngle float32
}

func (bullet *Bullet) Update(delta float32) {
	bullet.boundingBox.X += bullet.velocity * delta * float32(math.Cos(DegreeToRadian(float64(bullet.rotationAngle))))
	bullet.boundingBox.Y += bullet.velocity * delta * float32(math.Sin(DegreeToRadian(float64(bullet.rotationAngle))))
}

/*func (bullet Bullet) Draw(window *pixelgl.Window) {
	matrix := pixel.IM
	matrix = matrix.Moved(bullet.position)
	// matrix = matrix.Scaled(bullet.position, 1.0) // no need to scale, when scale is 1
	matrix = matrix.Rotated(bullet.position, bullet.rotationAngle)
	bullet.bulletSprite.Draw(window, matrix)
}*/

func (bullet Bullet) OutOfWindow() bool {
	return !((bullet.boundingBox.X >= 0) &&
		(bullet.boundingBox.Y >= 0) &&
		((bullet.boundingBox.X + bullet.boundingBox.W) <= float32(SCREEN_WIDTH)) &&
		((bullet.boundingBox.Y + bullet.boundingBox.H) <= float32(SCREEN_HEIGHT)))
}

/*
Almost same as the Bullet struct, it is also inheriting this interface...
Even if I define it, I will not use it most probably

type Tank interface {
	// Init()
	Update()
	Draw(window *pixelgl.Window)
}*/

type EnemyTank struct {
	tankTexture   *sdl.Texture
	rotationAngle float32
	boundingBox   sdl.FRect
	alive         bool
	noUpdateTime  float32
	timer         time.Time
}

func NewEnemyTank(tankTexture *sdl.Texture, width int32, height int32, initialRotationAngle float32, alive bool, noUpdateTime float32) EnemyTank {
	return EnemyTank{
		tankTexture:   tankTexture,
		rotationAngle: initialRotationAngle,
		boundingBox: sdl.FRect{
			X: 0.0,
			Y: 0.0,
			W: float32(width),
			H: float32(height),
		},
		alive:        alive,
		noUpdateTime: noUpdateTime,
		timer:        time.Now(),
	}
}

func SetPositions(enemyTanks []EnemyTank, playerTankBoundingBox sdl.FRect, r *rand.Rand) {
	for index, _ := range enemyTanks {
		if enemyTanks[index].alive {
			enemyTanks[index].boundingBox = GetPositionOfOneEnemyTank(enemyTanks[index].boundingBox, enemyTanks[:index], playerTankBoundingBox, r)
		}
	}
}

func GetPositionOfOneEnemyTank(enemyTankBoundingBox sdl.FRect, otherEnemyTanks []EnemyTank, playerTankBoundingBox sdl.FRect, r *rand.Rand) sdl.FRect {
	experimentalTankBoundingBox := sdl.FRect{
		X: r.Float32() * float32(SCREEN_WIDTH),
		Y: r.Float32() * float32(SCREEN_HEIGHT),
		W: enemyTankBoundingBox.W,
		H: enemyTankBoundingBox.H,
	}
	if !ValidPosition(experimentalTankBoundingBox, otherEnemyTanks, playerTankBoundingBox) {
		return GetPositionOfOneEnemyTank(enemyTankBoundingBox, otherEnemyTanks, playerTankBoundingBox, r)
	}
	return experimentalTankBoundingBox
}

func ValidPosition(experimentalTankBoundingBox sdl.FRect, otherEnemyTanks []EnemyTank, playerTankBoundingBox sdl.FRect) bool {
	for idx, _ := range otherEnemyTanks {
		if otherEnemyTanks[idx].alive &&
			(otherEnemyTanks[idx].boundingBox.HasIntersection(&experimentalTankBoundingBox) ||
				experimentalTankBoundingBox.HasIntersection(&playerTankBoundingBox) ||
				!((experimentalTankBoundingBox.X > 0.0) &&
					(experimentalTankBoundingBox.Y > 0.0) &&
					((experimentalTankBoundingBox.X + experimentalTankBoundingBox.W) < float32(SCREEN_WIDTH)) &&
					((experimentalTankBoundingBox.Y + experimentalTankBoundingBox.H) < float32(SCREEN_HEIGHT)))) {
			return false
		}
	}
	return true
}

/*func (tank EnemyTank) Update(delta float64, r *rand.Rand, playerTankPosition pixel.Vec) (EnemyTank, Bullet) {
	var bullet Bullet

	return tank, bullet
}*/

func (tank EnemyTank) MoveInRandomDir(delta float32, r *rand.Rand) EnemyTank {
	switch r.Intn(4) {
	case 0:
		tank.boundingBox.Y += LEVEL_0_ENEMY_TANK_VELOCITY * delta // DOWN
	case 1:
		tank.boundingBox.Y -= LEVEL_0_ENEMY_TANK_VELOCITY * delta // UP
	case 2:
		tank.boundingBox.X += LEVEL_0_ENEMY_TANK_VELOCITY * delta // RIGHT
	case 3:
		tank.boundingBox.X -= LEVEL_0_ENEMY_TANK_VELOCITY * delta // LEFT
	}
	return tank
}

func (tank EnemyTank) SpinAndShoot(delta float32, r *rand.Rand, playerTankPosition sdl.FPoint, bulletTexture *sdl.Texture, bulletWidth int32, bulletHeight int32) (EnemyTank, Bullet) {
	/*
		TODO :
		switch r.Intn(2) {
		case 0:
			displacementVector := playerTankPosition.Sub(tank.position) // SHOOT THE PLAYER
			tank.rotationAngle = displacementVector.Angle()
		case 1:*/
	tank.rotationAngle = r.Float32() * 360 /*(2.0 * math.Pi)*/ // SHOOT ANYWHERE RANDOMLY
	/*}*/
	return tank, Bullet{
		bulletTexture: bulletTexture,
		velocity:      BULLET_VELOCITY,
		boundingBox: sdl.FRect{
			X: tank.boundingBox.X,
			Y: tank.boundingBox.Y,
			W: float32(bulletWidth),
			H: float32(bulletHeight),
		},
		rotationAngle: tank.rotationAngle,
	}
}

/*This function seems to be very innocent, not mutating the receiver.
Actually, it changes the timer of the receiver. Be careful...*/
func (tank *EnemyTank) WillUpdate() bool {
	if time.Since(tank.timer).Seconds() >= float64(tank.noUpdateTime) {
		tank.timer = time.Now()
		return true
	}
	return false
}

func (tank *EnemyTank) Die(delta float32) {
	tank.alive = false
	tank.tankDieAnimation(delta)
}

func (tank EnemyTank) tankDieAnimation(delta float32) {
	// TODO
}

/*func (tank EnemyTank) Draw(window *sdl.Window) {
	if tank.alive {
		matrix := pixel.IM
		matrix = matrix.Moved(tank.position)
		// matrix = matrix.Scaled(tank.position, 1.0) // no need to scale, when scale is 1
		matrix = matrix.Rotated(tank.position, tank.rotationAngle)
		tank.tankSprite.Draw(window, matrix)
	}
}*/

type PlayerTank struct {
	tankTexture   *sdl.Texture
	rotationAngle float32
	boundingBox   sdl.FRect
}

/*func (tank PlayerTank) Update() PlayerTank {
	return nil
}*/

/*func (tank PlayerTank) Draw(window *sdl.Window) {
	matrix := pixel.IM
	matrix = matrix.Moved(tank.position)
	// matrix = matrix.Scaled(tank.position, 1.0) // no need to scale, when scale is 1
	matrix = matrix.Rotated(tank.position, tank.rotationAngle)
	tank.tankSprite.Draw(window, matrix)
}*/

func (tank *PlayerTank) Shoot(bulletTexture *sdl.Texture, bulletWidth int32, bulletHeight int32) Bullet {
	return Bullet{
		bulletTexture: bulletTexture,
		velocity:      BULLET_VELOCITY,
		boundingBox: sdl.FRect{
			X: (tank.boundingBox.X + tank.boundingBox.X + tank.boundingBox.W) / 2.0, // TODO: NOT WORKING, shooting from the centre of tank
			Y: (tank.boundingBox.Y + tank.boundingBox.Y + tank.boundingBox.H) / 2.0, // TODO: NOT WORKING, shooting from the centre of tank
			W: float32(bulletWidth),
			H: float32(bulletHeight),
		},
		rotationAngle: tank.rotationAngle,
	}
}

func (tank *PlayerTank) RotateClockWise(delta float32) *PlayerTank {
	result := *tank                                     // making a copy
	result.rotationAngle += TANK_ROTATION_ANGLE * delta // mutating that copy
	return &result                                      // returning pointer to that copy
}

func (tank *PlayerTank) RotateAntiClockWise(delta float32) *PlayerTank {
	result := *tank                                     // making a copy
	result.rotationAngle -= TANK_ROTATION_ANGLE * delta // mutating that copy
	return &result                                      // returning pointer to that copy
}

func (tank *PlayerTank) MoveUp(delta float32) *PlayerTank {
	result := *tank                                      // making a copy
	result.boundingBox.Y -= PLAYER_TANK_VELOCITY * delta // mutating that copy
	return &result                                       // returning pointer to that copy
}

func (tank *PlayerTank) MoveDown(delta float32) *PlayerTank {
	result := *tank                                      // making a copy
	result.boundingBox.Y += PLAYER_TANK_VELOCITY * delta // mutating that copy
	return &result                                       // returning pointer to that copy
}

func (tank *PlayerTank) MoveLeft(delta float32) *PlayerTank {
	result := *tank                                      // making a copy
	result.boundingBox.X -= PLAYER_TANK_VELOCITY * delta // mutating that copy
	return &result                                       // returning pointer to that copycopy
}

func (tank *PlayerTank) MoveRight(delta float32) *PlayerTank {
	result := *tank                                      // making a copy
	result.boundingBox.X += PLAYER_TANK_VELOCITY * delta // mutating that copy
	return &result                                       // returning pointer to that copy                                                                // returning pointer to that copy
}
