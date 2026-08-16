# GopherCon UK 2026: Vibe Code a 2D Game with Go and Gemini

![GopherCon UK 2026 Logo](logo.png)

> **Note**: The agent skills in this repository will not be updated here. Updated and maintained versions of these skills are available at [github.com/danicat/skills](https://github.com/danicat/skills).

Materials, agent skills, and tooling for the "Vibe Code a 2D Game with Go and Gemini" workshop at GopherCon UK 2026.

- 📑 **Slide Deck**: [Vibe Code a 2D Game using Go and Gemini](https://speakerdeck.com/danicat/vibe-code-a-2d-game-using-go-and-gemini)

## Session Overview

Participants build 2D games in Go using Ebitengine v2, integrated with Google Gemini models for image generation (Nano Banana), audio generation (Lyria 3), and automated Go quality tooling (GoDoctor).

### Design Approach

The skills and reference modules here are generic building blocks for any 2D game genre—arcade, puzzle, platformer, strategy, RPG, or rhythm games. Snippets like pathfinding, collision, or UI components are modular references that developers and coding agents adapt to fit specific mechanics.

## GopherCon UK Mini Game Jam

During the workshop session, participants competed in the **GopherCon UK Mini Game Jam**.

🏆 **Winner**: [**London Eco-Rider: Ice Cream & Bottle Patrol**](game-jam/london-eco-rider) by **Ivan S** ([@inesusvet](https://github.com/inesusvet)).

For complete submission details, see the [**Game Jam Showcase**](game-jam/README.md).

## Repository Structure

```text
.
├── .agents/
│   └── skills/                  # Agent skills for game development
│       ├── vibe-game-developer/ # Request router and orchestrator
│       ├── ebitengineer/        # Ebitengine v2 architecture and reference modules
│       ├── game-design/         # GDD generation and interactive game design
│       ├── godoctor/            # Go linting, formatting, testing, and mutation testing
│       ├── lyria/               # Music generation using Lyria 3 models
│       ├── nano-banana/         # Image generation and editing via Nano Banana
│       ├── procedural-art/      # Pure-code 2D graphics, vector math, and shaders
│       ├── procedural-composer/ # Pure-code audio synthesis and chiptune sound engine
│       ├── sprite-animation/    # Sprite sheet slicing, animation states, and Aseprite format
│       └── swarm-coding/        # Multi-agent task parallelization
├── game-jam/                    # GopherCon UK Mini Game Jam entries
│   └── london-eco-rider/        # Jam Winner: London Eco-Rider: Ice Cream & Bottle Patrol
├── check_env.sh                 # Pre-flight environment check script
└── README.md
```

## Agent Skills

### Orchestration & Design

- [`vibe-game-developer`](.agents/skills/vibe-game-developer/SKILL.md)  
  Routes user requests to the appropriate specialized skill based on the task (architecture, pure-code graphics, AI media generation, testing, or deployment).

- [`game-design`](.agents/skills/game-design/SKILL.md)  
  Guides game concept interviews via `/grill-me` and generates a structured Game Design Document ([`GDD.md`](.agents/skills/game-design/references/gdd_template.md)).

### Engine Architecture & Pure-Code Assets

- [`ebitengineer`](.agents/skills/ebitengineer/SKILL.md)  
  Core engineering rules for Ebitengine v2 games: 16:9 pixel-scaling canvas, 60 FPS delta-time loop, scene state machines, WebAssembly builds, and asset validation.  
  Includes reference guides for [project structure](.agents/skills/ebitengineer/references/project_structure.md), [server architecture](.agents/skills/ebitengineer/references/server_architecture.md), [physics & collisions](.agents/skills/ebitengineer/references/physics_and_collision.md), [tilemaps](.agents/skills/ebitengineer/references/tilemaps_and_levels.md), [UI layout](.agents/skills/ebitengineer/references/ui_and_hud.md), [input mapping](.agents/skills/ebitengineer/references/input_action_mapping.md), [entity pooling](.agents/skills/ebitengineer/references/entity_management.md), and [A* pathfinding](.agents/skills/ebitengineer/references/pathfinding_and_ai.md).

- [`procedural-art`](.agents/skills/procedural-art/SKILL.md)  
  Code-driven 2D drawing, matrix transformations (`GeoM`), color ramps, particle systems, and Kage shaders.

- [`sprite-animation`](.agents/skills/sprite-animation/SKILL.md)  
  Sprite sheet validation, grid slicing, animation controllers, Aseprite binary format parsing, and integration with `SolarLune/goaseprite`. Includes Go implementation in [`references/animation_controller.go`](.agents/skills/sprite-animation/references/animation_controller.go) and format reference in [`references/aseprite_format.md`](.agents/skills/sprite-animation/references/aseprite_format.md).

- [`procedural-composer`](.agents/skills/procedural-composer/SKILL.md)  
  Pure-code DSP synthesis engine, FM synthesis, ADSR envelopes, JSON sound specs, and CLI audio player.

### Generative AI Media Skills

- [`nano-banana`](.agents/skills/nano-banana/SKILL.md)  
  Image generation using Nano Banana models (`gemini-3.1-flash-lite-image`, `gemini-3.1-flash-image`, `gemini-3-pro-image`, `gemini-2.5-flash-image`). Handles pixel art sprites and character-consistent generation. Run via `uv run .agents/skills/nano-banana/scripts/banana.py`.

- [`lyria`](.agents/skills/lyria/SKILL.md)  
  44.1 kHz stereo audio generation using Lyria 3 models (`lyria-3-clip-preview` for 30-second loops and `lyria-3-pro-preview` for full tracks). Run via `uv run .agents/skills/lyria/scripts/lyria.py`.

### Quality & Workflow

- [`godoctor`](.agents/skills/godoctor/SKILL.md)  
  Go code quality enforcement using Google Go Style guidelines, SQL test analysis, and Selene mutation testing.

- [`swarm-coding`](.agents/skills/swarm-coding/SKILL.md)  
  Decomposes complex engineering features into isolated subtasks for parallel execution.

## Quick Start

### Requirements
- Go 1.26 or newer
- `uv` Python package runner
- Google Cloud SDK (`gcloud`) with Application Default Credentials configured (`gcloud auth application-default login`)

### Environment Pre-Flight Check
Run the included doctor script to verify all local tools before starting:

```bash
./check_env.sh
```

### CLI Tool Verification
Test the image and music generation CLI utilities:

```bash
# Verify Nano Banana image CLI
uv run .agents/skills/nano-banana/scripts/banana.py --help

# Verify Lyria music CLI
uv run .agents/skills/lyria/scripts/lyria.py --help
```