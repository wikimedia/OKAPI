import React, {useCallback, useState} from 'react'
import { TextField } from '@material-ui/core'
import Head from 'next/head'
import Header from '../components/Header'
import Footer from '../components/Footer'
import api from '../utils/api'
import validate from '../utils/validate'
import { ParamsType, ParamFields, ErrorType } from '../types/form'

const HomePage = () => {
  const [formSubmitted, setFormSubmitted] = useState(false);
  const [submitInProcess, setSubmitInProcess] = useState(false);
  const [isFirstTime, setIsFirstTime] = useState(true);
  const [errors, setErrors] = useState<ErrorType>({});
  const [fields, setFields] = useState<ParamsType>({
    email: null,
    fullName: null,
    company: null,
  });

  const fieldChange = useCallback((e: React.ChangeEvent<HTMLInputElement>) => {
    const { target: { id, value } } = e;

    if (!isFirstTime) {
      const formErrors = validate(fields);
      setErrors(formErrors)
    }

    setFields({...fields, [id]: value})
  }, [setFields, fields])

  const formSubmit = useCallback(async () => {
    const formErrors = validate(fields);

    if (!Object.keys(formErrors).length) {
      try {
        setSubmitInProcess(true)

        await api.post('/v1/leads/', fields)

        setFormSubmitted(true)
        setSubmitInProcess(false)
        setFields({
          email: null,
          fullName: null,
          company: null,
        })
        setErrors({})
      } catch(e) {
        setFormSubmitted(false)
        setSubmitInProcess(false)
        console.log(e);
      }
    } else {
      setErrors(formErrors)
      setIsFirstTime(false)
    }
  }, [fields]);

  const renderContent = () => {
    if (!formSubmitted) {
      return (
        <div className="form__wrapper">
          {submitInProcess && <div className="form__overlay" />}
          <div className="form__description">
            We will be launching the product later in 2021. If you or your company is interested in learning more, sign up for updates below:
          </div>
          <form id="lead">
            <div className="form__item">
              <TextField
                required
                fullWidth
                id={ParamFields.fullName}
                label="Full Name"
                onChange={(e: React.ChangeEvent<HTMLInputElement>) => fieldChange(e)}
                error={errors[ParamFields.fullName]}
                helperText={errors[ParamFields.fullName] ? 'Required' : ''}
                inputProps={{ maxLength: 250 }}
              />
            </div>
            <div className="form__item">
              <TextField
                required
                fullWidth
                id={ParamFields.company}
                label="Company"
                onChange={(e: React.ChangeEvent<HTMLInputElement>) => fieldChange(e)}
                error={errors[ParamFields.company]}
                helperText={errors[ParamFields.company] ? 'Required' : ''}
                inputProps={{ maxLength: 250 }}
              />
            </div>
            <div className="form__item">
              <TextField
                required
                fullWidth
                id={ParamFields.email}
                label="Work Email"
                onChange={(e: React.ChangeEvent<HTMLInputElement>) => fieldChange(e)}
                error={errors[ParamFields.email]}
                helperText={errors[ParamFields.email] ? 'Incorrect email' : ''}
                inputProps={{ maxLength: 250 }}
              />
            </div>
          </form>
          {renderButtons()}
          <br/>
          <p>If you’re a member of the press and have questions about Wikimedia Enterprise, please reach out to <a href="mailto:press@wikimedia.org">press@wikimedia.org</a>.</p>
        </div>
      );
    }

    return (
      <div className="form__description">
        <h3 className="h3">Thank you for your interest in Wikimedia Enterprise</h3>
        {/* Thank you for your interest in Wikimedia Enterprise. */}
        Our team will be in touch with you shortly.
      </div>
    );
  }

  const renderButtons = () => {
    if (!formSubmitted) {
      return (
        <div className="buttons">
          <button className="button" onClick={formSubmit} disabled={submitInProcess}>Submit</button>
        </div>
      )
    }

    return (
      <div className="buttons">
        <button className="button">Done</button>
      </div>
    )
  }

  return (
    <div className="layout">
      <Head>
        <title>Wikimedia Enterprise</title>
        <meta name="viewport" content="initial-scale=1.0, width=device-width" />
        <meta name="description" content="Enterprise-grade Wikimedia APIs, service, and support" />
        <meta name="keywords" content="Enterprise-grade Wikimedia APIs, service, and support" />
      </Head>
      <Header />
      <main className="main">
        <div className="wrapper">
          <h2 className="h2">
            Wikimedia Enterprise is a new product from the Wikimedia Foundation, the nonprofit that operates Wikipedia and other Wikimedia projects. Wikimedia Enterprise provides paid developer tools and services that make it easier for companies and organizations to consume and re-use Wikimedia data.
          </h2>
          {renderContent()}
        </div>
      </main>
      <Footer />
    </div>
  )
}

export default HomePage;
