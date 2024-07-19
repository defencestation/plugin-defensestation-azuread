package main

import (
    sync "sync"
    time "time"

    plugin "github.com/defensestation/pluginutils"
    azidentity "github.com/Azure/azure-sdk-for-go/sdk/azidentity"
    client "github.com/defensestation/azurehound/client"
    models "github.com/microsoftgraph/msgraph-sdk-go/models"
)

type AzureADPlugin struct {
    Options  *Options
    Credential *azidentity.ClientSecretCredential
    Client *client.AzureClient
    Plugin   *plugin.Plugin
    wg       sync.WaitGroup
}

type Options struct {
    SubscriptionID  string `json:"subscription_id"`
    TenantID  string `json:"tenant_id"`
    ClientID   string `json:"client_id"`
    ClientSecret   string `json:"client_secret"`
    UserEmailArray []string `json:"user_email_array"`
}

// types needed for msgraphsdk to put user, group, role, adminunit objects
type UserProperties struct {
    AboutMe *string
    AccountEnabled *bool
    Activities []models.UserActivityable
    AgeGroup *string
    AgreementAcceptances []models.AgreementAcceptanceable
    AppRoleAssignments []models.AppRoleAssignmentable
    AssignedLicenses []models.AssignedLicenseable
    AssignedPlans []models.AssignedPlanable
    Authentication models.Authenticationable
    AuthorizationInfo models.AuthorizationInfoable
    Birthday *time.Time
    BusinessPhones []string
    Calendar models.Calendarable
    CalendarGroups []models.CalendarGroupable
    Calendars []models.Calendarable
    CalendarView []models.Eventable
    Chats []models.Chatable
    City *string
    CloudClipboard models.CloudClipboardRootable
    CompanyName *string
    ConsentProvidedForMinor *string
    ContactFolders []models.ContactFolderable
    Contacts []models.Contactable
    Country *string
    CreatedDateTime *time.Time
    CreatedObjects []models.DirectoryObjectable
    CreationType *string
    CustomSecurityAttributes models.CustomSecurityAttributeValueable
    Department *string
    DeviceEnrollmentLimit *int32
    DeviceManagementTroubleshootingEvents []models.DeviceManagementTroubleshootingEventable
    DirectReports []models.DirectoryObjectable
    DisplayName *string
    Drive models.Driveable
    Drives []models.Driveable
    EmployeeExperience models.EmployeeExperienceUserable
    EmployeeHireDate *time.Time
    EmployeeId *string
    EmployeeLeaveDateTime *time.Time
    EmployeeOrgData models.EmployeeOrgDataable
    EmployeeType *string
    Events []models.Eventable
    Extensions []models.Extensionable
    ExternalUserState *string
    ExternalUserStateChangeDateTime *time.Time
    FaxNumber *string
    FollowedSites []models.Siteable
    GivenName *string
    HireDate *time.Time
    Identities []models.ObjectIdentityable
    ImAddresses []string
    InferenceClassification models.InferenceClassificationable
    Insights models.OfficeGraphInsightsable
    Interests []string
    IsResourceAccount *bool
    JobTitle *string
    JoinedTeams []models.Teamable
    LastPasswordChangeDateTime *time.Time
    LegalAgeGroupClassification *string
    LicenseAssignmentStates []models.LicenseAssignmentStateable
    LicenseDetails []models.LicenseDetailsable
    Mail *string
    MailboxSettings models.MailboxSettingsable
    MailFolders []models.MailFolderable
    MailNickname *string
    ManagedAppRegistrations []models.ManagedAppRegistrationable
    ManagedDevices []models.ManagedDeviceable
    Manager models.DirectoryObjectable
    MemberOf []models.DirectoryObjectable
    Messages []models.Messageable
    MobilePhone *string
    MySite *string
    Oauth2PermissionGrants []models.OAuth2PermissionGrantable
    OfficeLocation *string
    Onenote models.Onenoteable
    OnlineMeetings []models.OnlineMeetingable
    OnPremisesDistinguishedName *string
    OnPremisesDomainName *string
    OnPremisesExtensionAttributes models.OnPremisesExtensionAttributesable
    OnPremisesImmutableId *string
    OnPremisesLastSyncDateTime *time.Time
    OnPremisesProvisioningErrors []models.OnPremisesProvisioningErrorable
    OnPremisesSamAccountName *string
    OnPremisesSecurityIdentifier *string
    OnPremisesSyncEnabled *bool
    OnPremisesUserPrincipalName *string
    OtherMails []string
    Outlook models.OutlookUserable
    OwnedDevices []models.DirectoryObjectable
    OwnedObjects []models.DirectoryObjectable
    PasswordPolicies *string
    PasswordProfile models.PasswordProfileable
    PastProjects []string
    People []models.Personable
    PermissionGrants []models.ResourceSpecificPermissionGrantable
    Photo models.ProfilePhotoable
    Photos []models.ProfilePhotoable
    Planner models.PlannerUserable
    PostalCode *string
    PreferredDataLocation *string
    PreferredLanguage *string
    PreferredName *string
    Presence models.Presenceable
    Print models.UserPrintable
    ProvisionedPlans []models.ProvisionedPlanable
    ProxyAddresses []string
    RegisteredDevices []models.DirectoryObjectable
    Responsibilities []string
    Schools []string
    ScopedRoleMemberOf []models.ScopedRoleMembershipable
    SecurityIdentifier *string
    ServiceProvisioningErrors []models.ServiceProvisioningErrorable
    Settings models.UserSettingsable
    ShowInAddressList *bool
    SignInActivity models.SignInActivityable
    SignInSessionsValidFromDateTime *time.Time
    Skills []string
    Sponsors []models.DirectoryObjectable
    State *string
    StreetAddress *string
    Surname *string
    Teamwork models.UserTeamworkable
    Todo models.Todoable
    TransitiveMemberOf []models.DirectoryObjectable
    UsageLocation *string
    UserPrincipalName *string
    UserType *string
}


type GroupProperties struct {
    AcceptedSenders []models.DirectoryObjectable
    AllowExternalSenders *bool
    AppRoleAssignments []models.AppRoleAssignmentable
    AssignedLabels []models.AssignedLabelable
    AssignedLicenses []models.AssignedLicenseable
    AutoSubscribeNewMembers *bool
    Calendar models.Calendarable
    CalendarView []models.Eventable
    Classification *string
    Conversations []models.Conversationable
    CreatedDateTime *time.Time
    CreatedOnBehalfOf models.DirectoryObjectable
    Description *string
    DisplayName *string
    Drive models.Driveable
    Drives []models.Driveable
    Events []models.Eventable
    ExpirationDateTime *time.Time
    Extensions []models.Extensionable
    GroupLifecyclePolicies []models.GroupLifecyclePolicyable
    GroupTypes []string
    HasMembersWithLicenseErrors *bool
    HideFromAddressLists *bool
    HideFromOutlookClients *bool
    IsArchived *bool
    IsAssignableToRole *bool
    IsSubscribedByMail *bool
    LicenseProcessingState models.LicenseProcessingStateable
    Mail *string
    MailEnabled *bool
    MailNickname *string
    MemberOf []models.DirectoryObjectable
    Members []models.DirectoryObjectable
    MembershipRule *string
    MembershipRuleProcessingState *string
    MembersWithLicenseErrors []models.DirectoryObjectable
    Onenote models.Onenoteable
    OnPremisesDomainName *string
    OnPremisesLastSyncDateTime *time.Time
    OnPremisesNetBiosName *string
    OnPremisesProvisioningErrors []models.OnPremisesProvisioningErrorable
    OnPremisesSamAccountName *string
    OnPremisesSecurityIdentifier *string
    OnPremisesSyncEnabled *bool
    Owners []models.DirectoryObjectable
    PermissionGrants []models.ResourceSpecificPermissionGrantable
    Photo models.ProfilePhotoable
    Photos []models.ProfilePhotoable
    Planner models.PlannerGroupable
    PreferredDataLocation *string
    PreferredLanguage *string
    ProxyAddresses []string
    RejectedSenders []models.DirectoryObjectable
    RenewedDateTime *time.Time
    SecurityEnabled *bool
    SecurityIdentifier *string
    ServiceProvisioningErrors []models.ServiceProvisioningErrorable
    Settings []models.GroupSettingable
    Sites []models.Siteable
    Team models.Teamable
    Theme *string
    Threads []models.ConversationThreadable
    TransitiveMemberOf []models.DirectoryObjectable
    TransitiveMembers []models.DirectoryObjectable
    UniqueName *string
    UnseenCount *int32
    Visibility *string
}

type DirectoryRoleProperties struct {
    Description *string
    DisplayName *string
    Members []models.DirectoryObjectable
    RoleTemplateId *string
    ScopedMembers []models.ScopedRoleMembershipable
}

type AdminUnitProperties struct {
    Description *string
    DisplayName *string 
    Extensions []models.Extensionable
    Members []models.DirectoryObjectable
    ScopedRoleMembers []models.ScopedRoleMembershipable
    Visibility *string
}