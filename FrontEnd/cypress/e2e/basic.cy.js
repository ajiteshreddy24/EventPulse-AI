describe('Basic App Test', () => {
    it('loads home page', () => {
      cy.visit('/')
      cy.contains('Home').should('exist')
    })
  })