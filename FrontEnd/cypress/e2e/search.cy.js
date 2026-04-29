describe("Search Events", () => {
  it("filters events by search input", () => {
    cy.intercept("GET", "/api/events", [
      {
        id: 1,
        title: "AI Club Meetup",
        description: "Discussing machine learning",
        location: "Innovation Lab",
        event_date: "2026-05-01T18:00:00Z",
        capacity: 50,
        rsvpCount: 10,
        waitlistCount: 0,
        userHasRSVP: false,
        userOnWaitlist: false,
      },
      {
        id: 2,
        title: "Music Night",
        description: "Live performances",
        location: "Student Center",
        event_date: "2026-05-02T20:00:00Z",
        capacity: 50,
        rsvpCount: 12,
        waitlistCount: 0,
        userHasRSVP: false,
        userOnWaitlist: false,
      },
    ]).as("getEvents")

    cy.visit("/events")
    cy.wait("@getEvents")

    cy.get("[data-cy=search-input]").type("AI Club")

    cy.contains("AI Club Meetup").should("exist")
    cy.contains("Music Night").should("not.exist")
  })
})
