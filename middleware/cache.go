package middleware

import "github.com/gofiber/fiber/v2"

// PublicCache marks GET responses as publicly cacheable by browsers and any CDN
// (e.g. Cloudflare) sitting in front of the API. These endpoints serve the same
// anonymous, admin-managed data to everyone, so edge caching offloads most
// traffic before it reaches the server.
//
// s-maxage governs shared (CDN) caches; max-age governs the browser;
// stale-while-revalidate lets a cache serve slightly stale data while it
// refreshes in the background, avoiding latency spikes at expiry.
func PublicCache() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if c.Method() == fiber.MethodGet {
			c.Set(fiber.HeaderCacheControl, "public, max-age=120, s-maxage=300, stale-while-revalidate=60")
		}
		return c.Next()
	}
}

// NoStore prevents browsers and shared caches from storing a response. Applied
// to authenticated / per-user / admin endpoints, whose responses must never be
// served to another user or cached at the edge.
func NoStore() fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Set(fiber.HeaderCacheControl, "private, no-store")
		return c.Next()
	}
}
