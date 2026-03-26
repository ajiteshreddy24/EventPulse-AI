describe('Delete Event', () => {
    it('deletes an event', () => {
  
      cy.request('POST', 'http://localhost:5174/api/events', {
        title: 'Delete Unique Event',
        description: 'Test',
        location: 'Campus',
        event_date: new Date().toISOString()
      })
  
      cy.visit('/events')
  
      cy.contains('Delete Unique Event')
  
      cy.get('[data-cy=delete-btn]').first().click()
  
      cy.wait(500)
      cy.reload()
  
      cy.contains('Delete Unique Event').should('not.exist')
    })
  })