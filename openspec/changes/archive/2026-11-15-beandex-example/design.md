# Design

## Context

The experiment asks whether a tight design system plus real reference screens
raise agent output quality. The example must therefore be opinionated and
visually consistent, not a feature menu. See `proposal.md` for motivation.

## Goals / Non-Goals

**Goals:**

- A brief a human would actually review, with realistic data volume.
- References that teach density, hierarchy and restraint.
- No character or colour outside the example design system.

**Non-Goals:**

- Running beandex. It is a planned product, not a dependency.
- More than one screen.
- Pixel-perfect reproductions of lazygit, k9s or btop.

## Decisions

- **Reconstructed captures rather than live `tmux` captures.** The tools named in
  the briefing are not installed and a fixed capture is easier to review. Each
  file states it is a reconstruction, and the note names the principle rather
  than the exact layout.
- **Four layout directions written into the brief.** A flat grouped table, a
  two-pane list plus detail, a compact one-line-per-repo view and one free
  choice. This makes the comparison about structure, not colour.
- **Realistic sample data.** Three groups and twelve repos with paths such as
  `~/code/nivis`, counts and ages. Real strings expose truncation problems that
  `repo-1` never will.
- **Sizes only in `design.yaml`.** The brief repeats them for the agent, but the
  validator reads YAML. One source of truth for enforcement.
- **References are plain text in the allowed character set.** A reference that
  cannot be typed into a mockup teaches the wrong lesson.

## Risks / Trade-offs

- [Reconstructed references may encode my taste rather than the real tools] ->
  notes state the principle and the brief allows a free fourth direction.
- [Twelve repos may not fill 100x30 with four key groups] -> the brief asks for
  about three groups and twelve repos and allows empty space; the 70 percent
  text-fill rule keeps air in the layout.
- [The example will drift from the design schema] -> the e2e suite validates it
  on every run.

## Migration Plan

Greenfield. Archive creates the `example-project` main spec.
