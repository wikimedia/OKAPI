import React from 'react';
import Link from 'next/link'

const Header: React.FC = () => (
  <header className="header">
    <Link href="/">
      <img src="/images/header__logo.svg" className="header__logo" alt="Wikimedia enterprise logo" />
    </Link>
  </header>
)

export default Header;
