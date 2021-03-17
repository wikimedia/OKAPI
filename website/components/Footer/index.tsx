import React from 'react'
import { useRouter } from 'next/router'
import classnames from 'classnames'
import Link from 'next/link'

const links = [
  {
    title: 'Privacy policy',
    url: '/privacy-policy.html',
    key: 'privacy-policy',
    external: false
  },
  // {
  //   title: 'Source code',
  //   url: 'https://github.com/wikimedia/OKAPI',
  //   key: 'source-code',
  //   external: true
  // },
]

const Footer: React.FC = () => {
  const { asPath: pathname } = useRouter()

  return (
    <footer className="footer">
      <div className="footer__wrapper wrapper">
        <a className="footer__url link" href="https://wikimediafoundation.org/">
          <img src="/images/footer__logo.svg" alt="A WIKIMEDIA FOUNDATION COMPANY LOGO" className="footer__image" />
          A WIKIMEDIA FOUNDATION COMPANY
        </a>
        <div className="footer__links">
          {
            links.map(({key, title, url, external}) => {
              const linkClasses = classnames('footer__link link', {
                'link--active': pathname === url
              })

              if (external) {
                return <a className={linkClasses} href={url} target="_blank" key={key}>{title}</a>
              }

              return (
                <Link href={url} key={key}>
                  <a className={linkClasses}>{title}</a>
                </Link>
              )
            })
          }
        </div>
      </div>
    </footer>
  )
}

export default Footer;
