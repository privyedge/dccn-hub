package templates

// RegistrationTemplate provides the email template when there is a
// new registration
const RegistrationTemplate = `
<p>
	<b>Activate Your Ankr Account</b>
</p>
<p>
	Welcome! Thanks for joining Ankr Network. To access Ankr's DCCN service, please activate your account by clicking here: 
	<a href="https://{{.AppDomain}}/account-verify?code={{.Code}}&email={{.Email}}">Activate Now</a>.
</p>
<p>
	If clicking the link above doesn't work, please copy and paste the below URL in a new browser window.<br />
	https://{{.AppDomain}}/account-verify?code={{.Code}}&email={{.Email}}
</p>
<p>
	Having issues with setting up your account? Contact 
	<a href="mailto:account@ankr.network?subject=Account%20Activation%20issue">Ankr Network Account Support</a>.
</p>
`

// ForgotPasswordTemplate provides the email template when a user forgets his/her
// password and choose to reset
const ForgotPasswordTemplate = `
<p>
	<b>Hello Ankr User</b>,
</p>
<p>
	You are receiving this message because you or somebody else has attempted to reset 
	your password on Ankr Portal. If you did not request a new password, please disregard this message.
</p>
<p>
	To reset your password, please click here: 
	<a href="https://{{.AppDomain}}/reset-verify?code={{.Code}}&email={{.Email}}">Reset Password</a>
</p>
<p>
	If clicking the link above doesn't work, please copy and paste the below URL in a 
	new browser window.<br />
	https://{{.AppDomain}}/reset-verify?code={{.Code}}&email={{.Email}}
</p>
<p>
	Thanks,<br />
	The Ankr Team
</p>
`

// ChangeEmailTemplate provides the email template when a user choose to
// change his/her email
const ChangeEmailTemplate = `
<p>
	<b>Hello Ankr User</b>,
</p>
<p>
	You are receiving this message because you requested a change to the email address on your Ankr account to {{.NewEmail}}.
</p>
<p>
	To confirm this change, please click <a href="https://{{.AppDomain}}/email-verify?new_email={{.NewEmailEncoded}}">here</a>. 
	Note that this link will expire in 24 hours.
</p>
<p>
	If clicking the link above doesn't work, please copy and paste the below URL in a new browser window.<br />
	https://{{.AppDomain}}/email-verify?new_email={{.NewEmailEncoded}}</p>
<p>
	Thanks,<br /> 
	The Ankr Team
</p>
<p>
	If you did not authorize this change, please contact 
	<a href="mailto:account@ankr.network?subject=Change%20Email%20issue">Ankr Network Account Support</a>.
</p>
<hr />
<p>
	You are receiving this email because you took an action on 
	<a href="https://www.ankr.com/">Ankr.com</a>.<br />
</p>
`
