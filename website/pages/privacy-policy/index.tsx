import React from 'react'
import Head from 'next/head'
import Header from '../../components/Header'
import Footer from '../../components/Footer'

const PrivacyPage = () => (
  <div className="layout">
    <Head>
      <title>Wikimedia Enterprise Privacy Policy</title>
      <meta name="viewport" content="initial-scale=1.0, width=device-width" />
      <meta name="description" content="Privacy Policy for Wikimedia Enterprise" />
      <meta name="keywords" content="Privacy Policy for Wikimedia Enterprise" />
    </Head>
    <Header />
    <main className="main main--gap">
      <div className="wrapper">
        <h1 className="h1">Privacy Policy</h1>
        <p className="paragraph">We at Wikimedia LLC (“Wikimedia”) know you care about how your personal information is used and shared, and we take your privacy seriously. Please read the following to learn more about our Privacy Policy. By using or accessing enterprise.wikimedia.com or Wikimedia Enterprise (“Service”) in any manner, you acknowledge that you accept the practices and policies outlined in this Privacy Policy, and you hereby consent that we will collect, use, and share your information in the following ways.</p>

        <h3 className="h3">What does this Privacy Policy cover?</h3>
        <p className="paragraph">This Privacy Policy covers our treatment of personally identifiable information (“Personal Information”) that we gather when you are accessing or using our Service, but not to the practices of companies we don’t own or control, or people that we don’t manage. We gather various types of Personal Information from our users, as explained in more detail below, and we use this Personal Information internally in connection with our Service, including to personalize, provide and improve our services, to contact you, and to analyze how you use the Service.</p>

        <p className="paragraph">We do not knowingly collect or solicit Personal Information from anyone under the age of 13. If you are under 13, please do not attempt to register for the Service or send any Personal Information about yourself to us. If we learn that we have collected Personal Information from a child under age 13, we will delete that Personal Information as quickly as possible. If you believe that a child under 13 may have provided us Personal Information, please contact us at <a href="mailto:privacy@wikimedia.org" target="_blank">privacy@wikimedia.org</a>.</p>

        <h3 className="h3">Will Wikimedia ever change this Privacy Policy?</h3>
        <p className="paragraph">We’re constantly trying to improve our Service, so we may need to change this Privacy Policy from time to time as well, but we will alert you to changes by placing a notice on <a href="//enterprise.wikimedia.com">enterprise.wikimedia.com</a>, by sending you an email, and/or by some other means. Please note that if you’ve opted not to receive legal notice emails from us (or you haven’t provided us with your email address), those legal notices will still govern your use of the Service, and you are still responsible for reading and understanding them. If you use the Service after any changes to the Privacy Policy have been posted, that means you agree to all of the changes. Use of Personal Information we collect now is subject to the Privacy Policy in effect at the time such Personal Information is used.</p>

        <h3 className="h3">What Personal Information does Wikimedia Collect?</h3>
        <p className="paragraph">We receive and store any Personal Information you knowingly provide to us. Certain Personal Information (such as account registration information) may be required to register with us or to take advantage of some of our features.</p>

        <p className="paragraph">We also collect Personal Information automatically about your use of the Service, such as metadata related to your use of the Service, your preferences, your configurations, and other information gathered about your use of the Service.</p>

        <p className="paragraph">We may communicate with you if you’ve provided us the means to do so. For example, if you’ve given us your email address, we may send you emails on behalf of Wikimedia, or email you about your use of the Service. If you do not want to receive communications from us, please indicate your preference by sending an email to <a href="mailto:privacy@wikimedia.org" target="_blank">privacy@wikimedia.org</a>.</p>

        <h3 className="h3">Will Wikimedia Share Any of the Personal Information it Receives?</h3>
        <p className="paragraph">In some circumstances, we employ other companies and people to perform tasks on our behalf and need to share your Personal Information with them to provide products or services to you; for example, we may use a third-party mail management service to send you emails on our behalf. Unless we tell you differently, our agents do not have any right to use the Personal Information we share with them beyond what is necessary to assist us.</p>

        <p className="paragraph">We reserve the right to access, read, preserve, and disclose any Personal Information that we reasonably believe is necessary to comply with law or court order; enforce or apply this Privacy Policy and other agreements; or protect the rights, property, or safety of Wikimedia, our employees, our users, or others.</p>

        <h3 className="h3">Is Personal Information about me secure?</h3>
        <p className="paragraph">We endeavor to protect the privacy of your Personal Information we hold in our records, but unfortunately, we cannot guarantee complete security. Unauthorized entry or use, hardware or software failure, and other factors, may compromise the security of user information at any time.</p>

        <h3 className="h3">What Personal Information can I access?</h3>
        <p className="paragraph">The Service currently does not include the ability to access your Personal Information. If you have any questions about the Personal Information we have on file about you, please contact us at <a href="mailto:privacy@wikimedia.org" target="_blank">privacy@wikimedia.org</a>.</p>

        <p className="paragraph">Under California Civil Code Sections 1798.83-1798.84, California residents are entitled to ask us for a notice identifying the categories of Personal Information which we share with our affiliates and/or third parties for marketing purposes, and providing contact information for such affiliates and/or third parties. If you are a California resident and would like a copy of this notice, please submit a written request to 1 Montgomery Street, Suite 1600, San Francisco, CA 94104.</p>

        <h3 className="h3">What choices do I have?</h3>
        <p className="paragraph">You can always opt not to disclose Personal Information to us, but keep in mind some Personal Information may be needed to register with us or to take advantage of some of our features.</p>

        <h3 className="h3">What if I have questions about this policy?</h3>
        <p className="paragraph">If you have any questions or concerns regarding this Privacy Policy, please send us a detailed message to <a href="mailto:privacy@wikimedia.org" target="_blank">privacy@wikimedia.org</a> and we will try to resolve your concerns.</p>
      </div>
    </main>
    <Footer />
  </div>
)

export default PrivacyPage;
