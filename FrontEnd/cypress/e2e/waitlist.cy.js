describe("Event Waitlist", () => {
  it("joins the waitlist when an event is full", () => {
    cy.intercept("GET", "/api/events", [
      {
        id: 1,
        title: "Full Capacity Event",
        description: "Seats are already filled",
        location: "Auditorium",
        event_date: "2026-05-04T19:00:00Z",
        capacity: 2,
        rsvpCount: 2,
        waitlistCount: 1,
        userHasRSVP: false,
        userOnWaitlist: false,
      },
    ]).as("getEvents")

    cy.intercept("POST", "/api/events/1/waitlist", {
      statusCode: 200,
      body: { message: "Added to waitlist" },
    }).as("joinWaitlist")

    cy.visit("/events")
    cy.wait("@getEvents")

    cy.get("[data-cy=waitlist-btn]").first().click()
    cy.wait("@joinWaitlist")
  })
})
