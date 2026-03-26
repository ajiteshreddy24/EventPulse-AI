describe('Navigation', () => {
    it('navigates correctly between pages', () => {
      cy.visit('/')
  
      cy.contains('Create Event').click()
      cy.url().should('include', '/create')
  
      cy.contains('Upcoming Events').click()
      cy.url().should('include', '/events')
    })
  })