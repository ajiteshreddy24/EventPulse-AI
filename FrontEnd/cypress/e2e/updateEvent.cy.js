describe('Update Event', () => {
    it('updates an event', () => {
  
      cy.request('POST', 'http://localhost:5174/api/events', {
        title: 'Old Cypress Title',
        description: 'Old desc',
        location: 'Old location',
        event_date: new Date().toISOString()
      })
  
      cy.visit('/events')
  
      cy.contains('Old Cypress Title')
  
      cy.get('[data-cy=edit-btn]').first().click()
  
      cy.get('[data-cy=title-input]').clear().type('Updated Cypress Title')
      cy.get('[data-cy=update-btn]').click()
  
      cy.url().should('include', '/events')
      cy.contains('Updated Cypress Title')
    })
  })