package api

import (
	"errors"
	"github.com/28Pollux28/webcross3d/internal/puzzle"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/session"
	"time"
)

type Handler struct {
	SessionStore *session.Store
}

func (h *Handler) GetApiHealthz(ctx *fiber.Ctx) error {
	return ctx.SendString("ok")
}

func (h *Handler) GetApiPuzzles(ctx *fiber.Ctx) error {
	puzzles := puzzle.GetPuzzles()
	var summaries []PuzzleSummary
	for _, p := range puzzles {
		summaries = append(summaries, PuzzleSummary{
			Id:     p.ID,
			Name:   p.Name,
			Author: p.Author,
		})
	}

	return ctx.JSON(summaries)
}

func (h *Handler) GetApiPuzzlesId(ctx *fiber.Ctx, id string) error {
	p, ok := puzzle.GetPuzzle(id)
	if !ok {
		return ctx.Status(fiber.StatusNotFound).SendString("not found")
	}

	return ctx.JSON(Puzzle{
		Id:     p.ID,
		Name:   "Intro Puzzle",
		Author: "Admin",
		SizeX:  p.SizeX,
		SizeY:  p.SizeY,
		SizeZ:  p.SizeZ,
		Clues:  []ClueSet{}, // TODO: compute clues
	})
}

func (h *Handler) PostApiPuzzlesIdStart(ctx *fiber.Ctx, id string) error {
	// Get the old session and destroy it
	old, err := h.SessionStore.Get(ctx)
	if err == nil {
		_ = old.Destroy()
	}
	sess, err := h.SessionStore.Get(ctx)
	if err != nil {
		return err
	}

	p, ok := puzzle.GetPuzzle(id)
	if !ok {
		return ctx.Status(fiber.StatusNotFound).SendString("puzzle not found")
	}

	sess.Set("puzzle_id", id)
	sess.Set("lives", p.Lives)
	sess.Set("start_time", time.Now().UnixMilli())
	err = sess.Save()
	if err != nil {
		log.Debugf("failed to save session: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to save session",
		})
	}

	return ctx.JSON(fiber.Map{
		"lives": p.Lives,
	})
}

func (h *Handler) PostApiPuzzlesIdActions(ctx *fiber.Ctx, id string) error {
	sess, err := h.SessionStore.Get(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "session error"})
	}

	puzzleID := sess.Get("puzzle_id")
	if puzzleID != id {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "session mismatch"})
	}

	var action VoxelAction
	if err := ctx.BodyParser(&action); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid payload"})
	}

	p, ok := puzzle.GetPuzzle(id)
	if !ok {
		return ctx.Status(fiber.StatusNotFound).SendString("puzzle not found")
	}

	// Session state
	lives := sess.Get("lives").(int)

	removed := make(map[[3]int]bool)
	if v := sess.Get("removed"); v != nil {
		removed = v.(map[[3]int]bool)
	}

	coord := [3]int{action.X, action.Y, action.Z}

	if removed[coord] {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "voxel already removed"})
	}

	correct, err := p.ValidateVoxel(action.X, action.Y, action.Z)
	if err != nil {
		if errors.Is(err, puzzle.ErrOutOfBounds) {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "coordinates out of bounds"})
		}
		if errors.Is(err, puzzle.ErrIncorrectVoxel) {
			lives--
			sess.Set("lives", lives)
			err := sess.Save()
			if err != nil {
				return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to save session"})
			}

			return ctx.JSON(ActionResult{
				Success:        false,
				RemainingLives: lives,
				Completed:      false,
			})
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal error"})
	}
	if !correct {
		lives--
		sess.Set("lives", lives)
		err := sess.Save()
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to save session"})
		}

		return ctx.JSON(ActionResult{
			Success:        false,
			RemainingLives: lives,
			Completed:      false,
		})
	}

	// If the voxel is correct, we proceed
	removed[coord] = true
	sess.Set("removed", removed)

	// Check if puzzle is now complete
	completed := p.IsComplete(removed)

	if completed {
		startTimeMillis, _ := sess.Get("start_time").(int64)
		startTime := time.UnixMilli(startTimeMillis)
		elapsed := int(time.Since(startTime).Seconds())

		return ctx.JSON(ActionResult{
			Success:        true,
			RemainingLives: lives,
			Completed:      true,
			CompletionTime: &elapsed,
		})
	}
	err = sess.Save()
	if err != nil {
		log.Debugf("failed to save session: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to save session"})
	}

	return ctx.JSON(ActionResult{
		Success:        true,
		RemainingLives: lives,
		Completed:      false,
	})
}
