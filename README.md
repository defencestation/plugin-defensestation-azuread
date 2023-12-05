### plugin-defensestation-azuread

To configure Azure Plugin, you will need to register an app in your azure AD, and give it users read permisions. And grant admin access to it.
Steps are as follows:

#### Register an App in Azure AD:

- Log in to the Azure portal.
- Navigate to Azure Active Directory.
- Select "App registrations" and click on "New registration."
- Fill in the required information, such as the name and redirect URI. Note down the "Application (client) ID" as this is the application ID needed for configuration.

#### Grant User Read Permissions:

- In the Azure portal, go to the newly registered app.
- Navigate to "API permissions" and click on "Add a permission."
- Choose the appropriate API (e.g., Microsoft Graph) and select the required permissions, such as "User.Read."(Make sure it is Application type, not the Delegated Type)
- Grant admin consent for the permissions.

#### Create a Client Secret:

- In the Azure portal, go to the registered app.
- Navigate to "Certificates & Secrets" and click on "New client secret."
- Provide a description, select an expiry period, and click "Add." Note down the generated client secret value immediately.

#### Copy and Paste Information:

- Once the app is registered and permissions are granted, copy the following information:
    - Application (client) ID: [Copy from the app registration page.]
    - Directory (tenant) ID: [Available on the app registration page.]
    - Client Secret: [Copy the value created in step 4.]