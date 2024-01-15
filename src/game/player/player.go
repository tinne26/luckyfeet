package player

import "image"
import "math"

import "github.com/hajimehoshi/ebiten/v2"

import "github.com/tinne26/luckyfeet/src/game/context"
import "github.com/tinne26/luckyfeet/src/game/material/in"
import "github.com/tinne26/luckyfeet/src/game/material/au"
import "github.com/tinne26/luckyfeet/src/game/components/tile"
import "github.com/tinne26/luckyfeet/src/game/components/tile/tcsts"
import "github.com/tinne26/luckyfeet/src/game/player/motion"
import "github.com/tinne26/luckyfeet/src/game/carrot"

type Player struct {
	state State
	anim *motion.Animation
	x, y float64
	dir in.Direction
	lastActiveLayer int
	
	drawOpts ebiten.DrawImageOptions
	
	vertSpeed float64
	jumpSpeedGainLeft float64
	jumpingTicks int
	jumpHoldStopTick int
	didTicTac bool
	ticksInExtraGravity int
}

func New(ctx *context.Context) *Player {
	return &Player{
		state: StFalling,
		anim: ctx.Animations.InAir,
		dir: in.DirRight,
		lastActiveLayer: tcsts.LayerMain,
		jumpHoldStopTick: 9999,
		didTicTac: true,
	}
}

func (self *Player) Respawn(ctx *context.Context, tilemap *tile.Map) {
	self.changeState(ctx, StIdle, ctx.Animations.Idle)
	self.x = float64(tilemap.StartCol)*20 + 2
	self.y = float64(tilemap.StartRow)*20 - CollisionHeight + 11
	self.vertSpeed = 0
	self.jumpSpeedGainLeft = 0
	self.jumpingTicks = 0
	self.didTicTac = false
	self.ticksInExtraGravity = 0
	self.lastActiveLayer = tcsts.LayerMain
	
	row, col := tilemap.StartRow, tilemap.StartCol
	_, hasFrontTile := tilemap.GetTileIDAt(row, col, tcsts.LayerFront)
	if hasFrontTile { self.lastActiveLayer = tcsts.LayerFront }
}

func (self *Player) Update(ctx *context.Context, carrots *carrot.Inventory, tilemap *tile.Map) error {
	dir := ctx.Input.HorzDir()
	if dir != in.DirNone { self.dir = dir }

	switch self.state {
	case StIdle:
		if dir == in.DirNone {
			if self.detectAndProcessFalling(ctx, carrots, tilemap) { break }
			slipX := self.detectSlip(ctx, carrots, tilemap)
			if slipX != self.x { self.slipTowards(ctx, carrots, tilemap, slipX) }
		} else {
			self.changeState(ctx, StRunning, ctx.Animations.Running)
			self.applyRunningMotion(ctx, carrots, tilemap) // includes falling/slip detection too
		}
		
		// jump triggering
		if self.state == StIdle && ctx.Input.Trigger(in.ActionJump) {
			ctx.Audio.PlaySFX(au.SfxJump)
			self.changeState(ctx, StJumpingHold, ctx.Animations.InAir)
		}
	case StRunning:
		if dir == in.DirNone {
			//ctx.Audio.PlaySFX(au.SfxLand)
			self.changeState(ctx, StIdle, ctx.Animations.Idle)
		} else {
			self.dir = dir
			self.applyRunningMotion(ctx, carrots, tilemap) // includes falling/slip detection too
		}
		
		// jump triggering
		if self.state == StRunning && ctx.Input.Trigger(in.ActionJump) {
			ctx.Audio.PlaySFX(au.SfxJump)
			self.changeState(ctx, StJumpingHold, ctx.Animations.InAir)
		}
	case StJumpingHold:
		self.applyJumpMotion(ctx, carrots, tilemap, dir) // includes some state changes (top collisions, natural fall)
		if self.state != StJumpingHold { break }
		if ctx.Input.Trigger(in.ActionJump) && self.canTicTac(ctx, carrots, tilemap) {
			ctx.Audio.PlaySFX(au.SfxTicTac)
			self.changeState(ctx, StTicTacHold, ctx.Animations.InAir)
		} else if !ctx.Input.Pressed(in.ActionJump) {
			self.jumpHoldStopTick = self.jumpingTicks
			self.changeState(ctx, StJumpingInertial, ctx.Animations.InAir)
		}
	case StJumpingInertial:
		self.applyJumpMotion(ctx, carrots, tilemap, dir) // includes some state changes (top collisions, natural fall)
		if self.state != StJumpingInertial { break }
		if ctx.Input.Trigger(in.ActionJump) && self.canTicTac(ctx, carrots, tilemap) {
			ctx.Audio.PlaySFX(au.SfxTicTac)
			self.changeState(ctx, StTicTacHold, ctx.Animations.InAir)
			break
		}
	case StTicTacHold:
		self.applyJumpMotion(ctx, carrots, tilemap, dir) // includes some state changes (top collisions, natural fall)
		if self.state != StTicTacHold { break }
		if !ctx.Input.Pressed(in.ActionJump) {
			self.changeState(ctx, StTicTacInertial, ctx.Animations.InAir)
		}
	case StTicTacInertial:
		self.applyJumpMotion(ctx, carrots, tilemap, dir) // includes some state changes (top collisions, natural fall)
	case StFalling:
		self.applyFallMotion(ctx, carrots, tilemap, dir) // includes falling/slip detection too
		if self.state != StFalling { break }
		if ctx.Input.Trigger(in.ActionJump) && self.canTicTac(ctx, carrots, tilemap) {
			ctx.Audio.PlaySFX(au.SfxTicTac)
			self.changeState(ctx, StTicTacHold, ctx.Animations.InAir)
			break
		}
	default:
		panic("unknown player state")
	}

	// update animation after state update (should feel more responsive here)
	self.anim.Update(ctx.Audio)

	return nil
}

func (self *Player) Draw(canvas *ebiten.Image, ctx *context.Context) {
	frame := self.anim.GetCurrentFrame()
	if self.dir == in.DirLeft {
		self.drawOpts.GeoM.Scale(-1, 1)
		self.drawOpts.GeoM.Translate(float64(frame.Bounds().Dx()), 0)
	}

	ix, iy := self.getXYi()
	self.drawOpts.GeoM.Translate(float64(ix - CollisionXOffset), float64(iy - CollisionYOffset))
	canvas.DrawImage(frame, &self.drawOpts)
	self.drawOpts.GeoM.Reset()
}

func (self *Player) GetLightCenterPoint() (x, y int) {
	ix, iy := self.getXYi()
	switch self.dir {
	case in.DirRight : return ix - CollisionXOffset + 7, iy - CollisionYOffset + 18
	case in.DirLeft  : return ix - CollisionXOffset + 6, iy - CollisionYOffset + 18
	default:
		panic("broken code")
	}
}

func (self *Player) GetSpecialRect() image.Rectangle {
	ix, iy := self.getXYi()
	switch self.dir {
	case in.DirRight : return image.Rect(ix + 1, iy + 3, ix + 6, iy + CollisionHeight - 6)
	case in.DirLeft  : return image.Rect(ix + 3, iy + 3, ix + 8, iy + CollisionHeight - 6)
	default:
		panic("broken code")
	}
}

func (self *Player) HasFallen() bool {
	return self.y > 360 + CollisionHeight + 16 + 120
}

func (self *Player) BehindMain()  bool { return self.lastActiveLayer == tcsts.LayerBack }
func (self *Player) InFrontMain() bool { return self.lastActiveLayer != tcsts.LayerBack }

// --- private methods ----

func (self *Player) ensureAnimSet(ctx *context.Context, anim *motion.Animation) {
	if self.anim != anim {
		self.anim = anim
		self.anim.Rewind(ctx.Audio)
	}
}

func (self *Player) changeState(ctx *context.Context, newState State, anim *motion.Animation) {
	if newState == StRunning && anim == ctx.Animations.Running {
		self.anim = anim
		self.anim.RewindToLoop(ctx.Audio)
	} else {
		self.ensureAnimSet(ctx, anim)
	}
	self.state = newState
	
	switch newState {
	case StJumpingHold:
		self.didTicTac = false
		self.jumpSpeedGainLeft = JumpInitialSpeed
		self.vertSpeed = 0
		self.vertSpeed = self.nextJumpSpeed()
		self.ticksInExtraGravity = 0
		self.jumpingTicks = 0
		self.jumpHoldStopTick = 9999
	case StFalling:
		if self.vertSpeed > 0 { // necessary for bonk cases
			self.vertSpeed = 0
		}
	case StTicTacHold:
		self.vertSpeed = JumpInitialSpeed*0.76 - math.Abs(self.vertSpeed)/8.0
		self.didTicTac = true
	case StIdle, StRunning:
		self.vertSpeed = 0
	}
}

// Automatically changes states to falling or idle if necessary.
const runSpeed = 1.5/2.0 // 1.32
const runFirstStepsSpeed = 0.4/2.0
func (self *Player) applyRunningMotion(ctx *context.Context, carrots *carrot.Inventory, tilemap *tile.Map) {
	if self.detectAndProcessFalling(ctx, carrots, tilemap) { return }
	
	var speed float64 = runSpeed
	if self.anim.InPreLoopPhase() {
		speed = runFirstStepsSpeed
	}

	switch self.dir {
	case in.DirLeft:
		slipX := self.detectSlip(ctx, carrots, tilemap)
		if slipX > self.x {
			self.slipTowards(ctx, carrots, tilemap, slipX)
			return
		}

		target := min(max(self.x - speed, 0), 640 - CollisionWidth)
		for self.x != target {
			nextX := max(math.Floor(self.x - 0.0001), target)
			if self.detectCollisionAtX(ctx, carrots, tilemap, nextX) {
				self.changeState(ctx, StIdle, ctx.Animations.Idle)
				break
			} else {
				self.x = nextX
			}
		}
	case in.DirRight:
		slipX := self.detectSlip(ctx, carrots, tilemap)
		if slipX < self.x {
			self.slipTowards(ctx, carrots, tilemap, slipX)
			return
		}

		target := min(max(self.x + speed, 0), 640 - CollisionWidth)
		for self.x != target {
			nextX := min(math.Ceil(self.x + 0.0001), target)
			if self.detectCollisionAtX(ctx, carrots, tilemap, nextX) {
				self.changeState(ctx, StIdle, ctx.Animations.Idle)
				break
			} else {
				self.x = nextX
			}
		}
	default:
		panic("broken code")
	}
}

// Automatically changes states for landing/running if necessary.
func (self *Player) applyFallMotion(ctx *context.Context, carrots *carrot.Inventory, tilemap *tile.Map, dir in.Direction) {
	// detect landing at current point for safety
	if self.detectAndProcessLandingAtY(ctx, carrots, tilemap, self.y + CollisionHeight, dir) {
		ctx.Audio.PlaySFX(au.SfxLand)
		return
	}
	
	// get target x and y coords
	targetY := self.y - self.nextFallSpeed()
	targetX := self.x
	horzSpeed := runSpeed + AirExtraHorzSpeed
	if self.didTicTac { horzSpeed += AirExtraHorzSpeed*TicTacAirHorzSpeedMult }
	switch dir {
	case in.DirLeft  : targetX -= horzSpeed
	case in.DirRight : targetX += horzSpeed
	}
	targetX = min(max(targetX, 0), 640 - CollisionWidth)

	// fall until reaching target or can't fall no more
	for self.y < targetY {
		// apply horz movement
		if dir != in.DirNone {
			newX := self.x
			switch dir {
			case in.DirLeft  : newX = max(math.Floor(self.x - 0.0001), targetX)
			case in.DirRight : newX = min(math.Ceil(self.x + 0.0001), targetX)
			}
			if self.detectCollisionAtX(ctx, carrots, tilemap, newX) {
				dir = in.DirNone
			} else {
				self.x = newX
			}
		}
		
		// apply vert movement
		self.y = min(math.Ceil(self.y + 0.0001), targetY)
		if self.detectAndProcessLandingAtY(ctx, carrots, tilemap, self.y + CollisionHeight, dir) {
			ctx.Audio.PlaySFX(au.SfxLand)
			break
		}
	}

	// apply remaining horz movement
	if dir != in.DirNone {
		for self.x != targetX {
			newX := self.x
			switch dir {
			case in.DirLeft  : newX = max(math.Floor(self.x - 0.0001), targetX)
			case in.DirRight : newX = min(math.Ceil(self.x + 0.0001), targetX)
			default: panic("broken code")
			}
			if self.detectCollisionAtX(ctx, carrots, tilemap, newX) { break }
			self.x = newX
		}
	}
}

func (self *Player) endFall(ctx *context.Context, dir in.Direction) {
	if dir == in.DirNone {
		self.changeState(ctx, StIdle, ctx.Animations.Idle)
	} else {
		self.changeState(ctx, StRunning, ctx.Animations.Running)
	}
}

// Automatically changes states to falling if necessary.
func (self *Player) applyJumpMotion(ctx *context.Context, carrots *carrot.Inventory, tilemap *tile.Map, dir in.Direction) {
	// get jump speed and start falling if jump speed is 0
	self.jumpingTicks += 1
	currentJumpSpeed := self.nextJumpSpeed()
	if currentJumpSpeed <= 0 {
		self.y -= currentJumpSpeed
		self.changeState(ctx, StFalling, ctx.Animations.InAir)
		return
	}

	// get target x and y coords
	targetY := self.y - currentJumpSpeed
	targetX := self.x
	horzSpeed := runSpeed + AirExtraHorzSpeed
	if self.didTicTac { horzSpeed += AirExtraHorzSpeed*TicTacAirHorzSpeedMult }
	switch dir {
	case in.DirLeft  : targetX -= horzSpeed
	case in.DirRight : targetX += horzSpeed
	}
	targetX = min(max(targetX, 0), 640 - CollisionWidth)

	// go up until reaching target or hitting something
	for self.y > targetY {
		// apply horz movement
		newX := self.x
		if dir != in.DirNone {
			switch dir {
			case in.DirLeft  : newX = max(math.Floor(self.x - 0.0001), targetX)
			case in.DirRight : newX = min(math.Ceil(self.x + 0.0001), targetX)
			}
		}
		
		// apply vert movement
		newY := max(math.Floor(self.y - 0.0001), targetY)
		if self.detectCollisionAt(ctx, carrots, tilemap, newX, newY) {
			self.changeState(ctx, StFalling, ctx.Animations.InAir)
			break
		}
		self.x, self.y = newX, newY
	}

	// apply remaining horz movement
	if dir != in.DirNone {
		for self.x != targetX {
			newX := self.x
			switch dir {
			case in.DirLeft  : newX = max(math.Floor(self.x - 0.0001), targetX)
			case in.DirRight : newX = min(math.Ceil(self.x + 0.0001), targetX)
			default: panic("broken code")
			}
			if self.detectCollisionAtX(ctx, carrots, tilemap, newX) { break }
			self.x = newX
		}
	}
}

const JumpInitialSpeed = 2.4
const DefaultGravity = 0.046
const ExtraGravity = 0.10
const MaxFallSpeed = JumpInitialSpeed*1.33
const AirExtraHorzSpeed = 0.2
const TicTacAirHorzSpeedMult = 0.36
func (self *Player) nextJumpSpeed() float64 {
	switch self.state {
	case StJumpingHold, StJumpingInertial:
		return self.nextNormalJumpSpeed()
	case StTicTacHold, StTicTacInertial:
		return self.nextTicTacJumpSpeed()
	default:
		panic("broken code")
	}
}

func (self *Player) nextNormalJumpSpeed() float64 {
	if self.jumpSpeedGainLeft > 0 {
		gain := self.jumpSpeedGainLeft*0.24
		self.vertSpeed += gain
		self.jumpSpeedGainLeft -= gain
	}
	self.vertSpeed -= DefaultGravity
	if self.jumpingTicks > self.jumpHoldStopTick {
		diff := self.jumpingTicks - self.jumpHoldStopTick
		self.ticksInExtraGravity = max(self.ticksInExtraGravity, diff)
		self.vertSpeed -= ExtraGravity
	}
	return self.vertSpeed
}

func (self *Player) nextTicTacJumpSpeed() float64 {
	self.vertSpeed -= DefaultGravity*0.76
	if self.state == StTicTacInertial {
		self.vertSpeed -= DefaultGravity
	}
	return self.vertSpeed
}

func (self *Player) nextFallSpeed() float64 {
	self.vertSpeed -= DefaultGravity
	if self.ticksInExtraGravity > 0 {
		self.ticksInExtraGravity -= 1
		self.vertSpeed -= ExtraGravity
	}
	self.vertSpeed = max(self.vertSpeed, -MaxFallSpeed)

	return self.vertSpeed
}

// Returns true if the player is starting to fall.
func (self *Player) detectAndProcessFalling(ctx *context.Context, carrots *carrot.Inventory, tilemap *tile.Map) bool {
	ox, fx, y := self.getLandingZone()
	var landed bool = tilemap.HasLandingFor(ctx, carrots, ox, fx, y, tcsts.LayerMain)
	if !landed && self.lastActiveLayer != tcsts.LayerBack {
		landed = tilemap.HasLandingFor(ctx, carrots, ox, fx, y, tcsts.LayerFront)
	}
	if landed { return false }
	self.changeState(ctx, StFalling, ctx.Animations.InAir)
	return true
}

func (self *Player) detectCollisionAt(ctx *context.Context, carrots *carrot.Inventory, tilemap *tile.Map, x, y float64) bool {
	rect := self.collisionRect()
	xshift := int(x) - rect.Min.X
	yshift := int(y) - rect.Min.Y
	rect = rect.Add(image.Pt(xshift, yshift))
	switch self.lastActiveLayer {
	case tcsts.LayerMain:
		return tilemap.Collides(ctx, carrots, rect, tcsts.LayerMain)
	case tcsts.LayerFront:
		return tilemap.Collides(ctx, carrots, rect, tcsts.LayerFront)
	default: // back layer
		return false
	}
}

func (self *Player) detectCollisionAtX(ctx *context.Context, carrots *carrot.Inventory, tilemap *tile.Map, x float64) bool {
	return self.detectCollisionAt(ctx, carrots, tilemap, x, self.y)
}

func (self *Player) detectAndProcessLandingAtY(ctx *context.Context, carrots *carrot.Inventory, tilemap *tile.Map, y float64, dir in.Direction) bool {
	ox, fx, _ := self.getLandingZone()
	if tilemap.HasLandingFor(ctx, carrots, ox, fx, int(y), tcsts.LayerMain) {
		self.endFall(ctx, dir)
		self.lastActiveLayer = tcsts.LayerMain
		return true
	} else if self.lastActiveLayer != tcsts.LayerBack && tilemap.HasLandingFor(ctx, carrots, ox, fx, int(y), tcsts.LayerFront) {
		self.endFall(ctx, dir)
		self.lastActiveLayer = tcsts.LayerFront
		return true
	}
	return false
}

// Returns the slip x, which will be == self.x if no slip is happening.
func (self *Player) detectSlip(ctx *context.Context, carrots *carrot.Inventory, tilemap *tile.Map) float64 {
	lox, lfx, rox, rfx, y := self.getFootLandingZones()
	
	const slipSpeed = 0.3
	switch self.dir {
	case in.DirRight:
		if !tilemap.HasLandingFor(ctx, carrots, rox, rfx, y, self.lastActiveLayer) {
			return min(max(self.x + slipSpeed, 0), 640 - CollisionWidth)
		} else if !tilemap.HasLandingFor(ctx, carrots, lox, lfx, y, self.lastActiveLayer) {
			return min(max(self.x - slipSpeed, 0), 640 - CollisionWidth)
		} else {
			return self.x
		}
	case in.DirLeft:
		if !tilemap.HasLandingFor(ctx, carrots, lox, lfx, y, self.lastActiveLayer) {
			return min(max(self.x - slipSpeed, 0), 640 - CollisionWidth)
		} else if !tilemap.HasLandingFor(ctx, carrots, rox, rfx, y, self.lastActiveLayer) {
			return min(max(self.x + slipSpeed, 0), 640 - CollisionWidth)
		} else {
			return self.x
		}
	default:
		panic("broken code")
	}
}

// Slips towards the given x, which must be at dist <= 1.0 from self.x.
// It automatically detects collisions to avoid slips if necessary.
func (self *Player) slipTowards(ctx *context.Context, carrots *carrot.Inventory, tilemap *tile.Map, slipX float64) {
	// safety assertions
	slipX = min(max(slipX, 0), 640 - CollisionWidth)
	diff := slipX - self.x
	if diff < 0 { diff = -diff }
	if diff > 1.0 { panic("precondition violation") }

	// only slip if we don't collide going towards slipX
	if !self.detectCollisionAtX(ctx, carrots, tilemap, slipX) {
		self.x = slipX
	}
}

const CollisionXOffset = 3
const CollisionYOffset = 7
const CollisionWidth   = 9
const CollisionHeight  = 28
func (self *Player) collisionRect() image.Rectangle {
	ix, iy := self.getXYi()
	return image.Rect(ix, iy, ix + CollisionWidth, iy + CollisionHeight)
}
func (self *Player) getLandingZone() (ox, fx, y int) {
	ix, iy := self.getXYi()
	switch self.dir {
	case in.DirRight : return ix + 0, ix + 7, iy + CollisionHeight
	case in.DirLeft  : return ix + 2, ix + 9, iy + CollisionHeight
	default:
		panic("broken code")
	}
	return ix, ix + CollisionWidth, iy + CollisionHeight
}

func (self *Player) getFootLandingZones() (lox, lfx, rox, rfx, y int) {
	ix, iy := self.getXYi()
	switch self.dir {
	case in.DirRight : return ix + 0, ix + 4, ix + 3, ix + 7, iy + CollisionHeight
	case in.DirLeft  : return ix + 2, ix + 6, ix + 5, ix + 9, iy + CollisionHeight
	default:
		panic("broken code")
	}
}

func (self *Player) getTicTacRect() image.Rectangle {
	ix, iy := self.getXYi()
	return image.Rect(ix + 2, iy + 20, ix + CollisionWidth - 2, iy + CollisionHeight - 3)
}

func (self *Player) getXYi() (int, int) {
	return int(self.x), int(self.y)
}

func (self *Player) canTicTac(ctx *context.Context, carrots *carrot.Inventory, tilemap *tile.Map) bool {
	if self.didTicTac { return false }
	if self.state == StFalling && self.vertSpeed < -JumpInitialSpeed { return false }

	rect := self.getTicTacRect()
	switch self.lastActiveLayer {
	case tcsts.LayerMain:
		if tilemap.Collides(ctx, carrots, rect, tcsts.LayerBack) {
			self.lastActiveLayer = tcsts.LayerBack
			return true
		}
	case tcsts.LayerFront:
		if tilemap.Collides(ctx, carrots, rect, tcsts.LayerMain) {
			//self.lastActiveLayer = tcsts.LayerMain
			return true
		}
	}

	return false
}
