import React from 'react';
import { MuiThemeProvider, createMuiTheme } from '@material-ui/core'
import CssBaseline from '@material-ui/core/CssBaseline'
import '../scss/index.scss'

const theme = createMuiTheme({
  overrides: {
    MuiFormLabel: {
      asterisk: {
        color: '#db3131',
        '&$error': {
          color: '#db3131'
        },
      }
    }
  }
})

interface Props {
  pageProps: any,
  Component: any
}

const WikimediaEnterprise: React.FC<Props> = ({ Component, pageProps }) => {

  React.useEffect(() => {
    const jssStyles = document.querySelector('#jss-server-side');
    document.body.style.visibility = 'visible';

    jssStyles?.parentElement?.removeChild(jssStyles);
  }, []);

  return (
    <MuiThemeProvider theme={theme}>
      <CssBaseline />
      <Component {...pageProps} />
    </MuiThemeProvider>
  )
}

export default WikimediaEnterprise;
