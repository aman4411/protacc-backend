package handler

import (
	"github.com/aman4411/protacc-backend/cache"
	"github.com/gofiber/fiber/v2"
)

// CacheHandler exposes read/purge operations on the in-process cache for the
// admin panel. Admin-only (mounted behind Protected + RequireRole).
type CacheHandler struct {
	cache *cache.Cache
}

func NewCacheHandler(c *cache.Cache) *CacheHandler {
	return &CacheHandler{cache: c}
}

// Stats returns the live entry count, total and per family, for the admin UI.
func (h *CacheHandler) Stats(c *fiber.Ctx) error {
	total, families := h.cache.Stats()
	return c.JSON(fiber.Map{"total": total, "families": families})
}

// Purge clears the whole cache, or a single family when ?prefix= is supplied
// (e.g. ?prefix=services:). Returns the number of entries removed.
func (h *CacheHandler) Purge(c *fiber.Ctx) error {
	prefix := c.Query("prefix")
	if prefix != "" {
		return c.JSON(fiber.Map{"purged": h.cache.InvalidatePrefix(prefix), "prefix": prefix})
	}
	return c.JSON(fiber.Map{"purged": h.cache.Clear()})
}
