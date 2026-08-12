# Game Design Document: London Eco-Rider: Ice Cream & Bottle Patrol

## 1. Elevator Pitch & Core Concept
**London Eco-Rider: Ice Cream & Bottle Patrol** is a fast-paced, Mario-inspired 2D side-scrolling platformer set against the vibrant backdrop of downtown London. Players control Leo, a energetic young boy on his trusty bicycle, riding past iconic landmarks like Big Ben, red double-decker buses, and classic red telephone boxes.

Leo's dual mission:
1. **Collect delicious Ice Cream Cones** for instant score bonuses and temporary Turbo Speed Boosts!
2. **Collect Plastic Bottles** to fill the Eco-Recycling Meter, cleaning up London's streets and unlocking high-score multipliers!

---

## 2. Gameplay Loop & Mechanics

### Core Action Cycle
- **Action**: Pedal forward, control momentum, and execute well-timed bunny hops / jumps off road ramps.
- **Challenge**: Navigate around traffic obstacles (black cabs, parked double-decker buses, construction cones, potholes, pesky London pigeons).
- **Reward**: Collect Ice Cream Cones for speed & points; collect Plastic Bottles to fill the Recycling Meter for bonus multipliers and Eco Hero status.

### Physics & Momentum
- **Bike Momentum**: Realistic bicycle acceleration, momentum retention, and smooth friction.
- **Bunny Hop Jump**: Variable jump height based on button hold duration. Slope momentum off ramps.
- **Recycling Meter**: Every plastic bottle collected fills the Eco Meter (0-100%). At 100%, triggers 10 seconds of **Eco-Fever Mode** (double points + invincibility shield).

### Win & Loss Conditions
- **Victory Condition**: Reach the end of the London downtown circuit before time runs out or collect 50 Plastic Bottles.
- **Loss Condition**: Colliding with major hazards (double-decker bus front, deep potholes) depletes energy. Game Over occurs when energy hits 0% or the 180-second clock expires.

---

## 3. Controls & Input Mapping

| Action | Keyboard Input | Gamepad Input | Virtual Touch / WASM |
| :--- | :--- | :--- | :--- |
| **Move Left / Right** | Left/Right Arrows or `A` / `D` | D-Pad Left / Right | On-screen Left/Right Buttons |
| **Bunny Hop / Jump** | Spacebar / Up Arrow or `W` | Action Button A / South | On-screen Jump Button |
| **Turbo Boost** | Left Shift or `J` | Action Button X / West | On-screen Boost Button |
| **Pause / Menu** | `Escape` or `P` | Start Button | Pause Icon |
| **Fullscreen** | `F11` or `Alt+Enter` | - | - |

---

## 4. Visual Art Strategy (Gemini AI Pixel Art + Ebitengine)
- **Main Character**: Leo the Little Boy riding a yellow & blue bicycle with a recycling basket.
- **London Landmarks**: Pixel art backdrops featuring Big Ben, Westminster Palace, Red Telephone Boxes, Tower Bridge, and Red Double-Decker Buses.
- **Items & Collectibles**:
  - 🍦 **Ice Cream Cones**: Creamy vanilla & strawberry cone sprite.
  - 🍾 **Plastic Bottles**: Clear green/blue eco-friendly recycling bottles.
- **Asset Pipeline**: Gemini AI Nano Banana model (`gemini-3.1-flash-image`) to generate pixel art sprite sheets and background tiles, embedded via `embed.FS`.

---

## 5. Audio & Soundscape Strategy
- **Sound Effects (SFX)**: Pure-code procedural audio synthesis (`procedural-composer`) for zero latency and crisp audio:
  - Bicycle Bell / Ring ring
  - Jump / Bunny Hop WHOOSH
  - Ice Cream Pick-Up Slurp/Crunch
  - Plastic Bottle Recycling Ding!
  - Turbo Boost Whir
  - Hazard Collision Bump
- **Background Music (BGM)**: Upbeat 6-channel chiptune theme capturing London's cheerful, energetic city vibe.

---

## 6. Game State Progression (FSM)

```text
[ Boot State ] --> [ Title Screen / Attract Mode ] --> [ Gameplay Level ] --> [ Stage Clear / Game Over ]
                             ^                                                        |
                             └────────────────────────────────────────────────────────┘
```

---

## 7. Technical Architecture (Ebitengine v2)
- **Language & Engine**: Go 1.22+ with `github.com/hajimehoshi/ebiten/v2`.
- **Target Resolution**: 16:9 Widescreen Virtual Canvas (`640x360` virtual pixels scaled smoothly to any window size).
- **Structure**:
  - `cmd/game/main.go`: Binary entry point & Ebitengine initialization.
  - `internal/game/`: Main game struct, FSM state controller.
  - `internal/entity/`: Player bike entity, collectibles (bottles, ice cream), obstacles, particle systems.
  - `internal/physics/`: Momentum physics, AABB collisions, ramp slopes.
  - `internal/render/`: Parallax background scrolling, HUD, sprite rendering.
  - `internal/audio/`: Procedural sound synthesizer and audio manager.
  - `assets/`: Embedded PNG sprite assets, fonts, sound configs.
