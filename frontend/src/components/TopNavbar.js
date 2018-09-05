import React from 'react';
import { Navbar, Nav, NavItem } from 'react-bootstrap/lib';

function TopNavbar() {
  return (
    <Navbar inverse fixedTop>
      <Navbar.Header>
        <Navbar.Brand>
          <a href="/">Livy</a>
        </Navbar.Brand>
        <Navbar.Toggle />
      </Navbar.Header>
      {/*<Navbar.Collapse>*/}
        <Nav pullRight>
          <NavItem href="https://github.com/lauperch/livy">GITHUB</NavItem>
        </Nav>
      {/**</Navbar.Collapse>*/}
    </Navbar>
  );
}

export default TopNavbar;
