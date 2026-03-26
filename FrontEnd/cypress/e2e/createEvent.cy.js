describe('Create Event', () => {
    it('creates a new event successfully', () => {
  
      cy.intercept('POST', '/api/events').as('createEvent')
  
      cy.visit('/create')
  
      cy.get('[data-cy=title-input]').type('Cypress Create Test')
      cy.get('[data-cy=description-input]').type('Testing create flow')
      cy.get('[data-cy=location-input]').type('Campus Hall')
      cy.get('[data-cy=date-input]').type('2026-04-01T18:00')
  
      cy.get('[data-cy=create-btn]').click()
  
      cy.wait('@createEvent')
  
      cy.url().should('include', '/events')
      cy.contains('Cypress Create Test')
    })
  })