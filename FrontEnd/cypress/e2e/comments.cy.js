describe("Event Comments", () => {
  it("loads and toggles comments for an event", () => {
    cy.intercept("GET", "/api/events", [
      {
        id: 1,
        title: "Hackathon",
        description: "Build projects together",
        location: "Tech Hub",
        event_date: "2026-05-03T16:00:00Z",
        capacity: 40,
        rsvpCount: 18,
        waitlistCount: 0,
        userHasRSVP: false,
        userOnWaitlist: false,
      },
    ]).as("getEvents")

    cy.intercept("GET", "/api/events/1/comments", [
      {
        id: 9,
        event_id: 1,
        user_id: 3,
        content: "Looking forward to this event",
        created_at: "2026-05-01T10:00:00Z",
        user: {
          id: 3,
          name: "Cypress User",
          email: "cypress@example.com",
        },
      },
    ]).as("getComments")

    cy.visit("/events")
    cy.wait("@getEvents")

    cy.get("[data-cy=comments-toggle]").first().click()
    cy.wait("@getComments")
    cy.contains("Looking forward to this event").should("exist")

    cy.get("[data-cy=comments-toggle]").first().click()
    cy.contains("Looking forward to this event").should("not.exist")
  })
})
