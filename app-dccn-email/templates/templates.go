package templates

// RegistrationTemplate provides the email template when there is a
// new registration
const RegistrationTemplate = `
<p>
	<b>Activate Your Ankr Account</b>
</p>
<p>
	Welcome! Thanks for joining Ankr Network. To access Ankr's DCCN service, please activate your account by clicking here: 
	<a href="https://app.ankr.network/account_verify?code={{.Code}}&id={{.ID}}">Activate Now</a>.
</p>
<p>
	If clicking the link above doesn't work, please copy and paste the below URL in a new browser window.<br />
	https://app.ankr.network/account_verify?code={{.Code}}&id={{.ID}}
</p>
<p>
	Having issues with setting up your account? Contact 
	<a href="mailto:account@ankr.network?subject=Account%20Activation%20issue">Ankr Network Account Support</a>.
</p>
`
