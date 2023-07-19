# Setting up an OAuth application in Azure

### Step 1: Create Load Test App in Azure

1. Sign into [Azure portal](https://portal.azure.com) using an admin Azure account.
2. Navigate to [App Registrations](https://portal.azure.com/#blade/Microsoft_AAD_IAM/ActiveDirectoryMenuBlade/RegisteredApps)
3. Click on **New registration** at the top of the page.

![image](https://user-images.githubusercontent.com/6913320/76347903-be67f580-62dd-11ea-829e-236dd45865a8.png)

4. Fill out the form with the following values:

- Name: `MS Teams Load Testing`
- Supported account types: Default value (Single tenant)

5. Navigate to **Certificates & secrets** in the left pane.

6. Click on **New client secret**. Then click on **Add**, and copy the new secret on the bottom right corner of the screen. We'll use this value later in the config file.

![image](https://user-images.githubusercontent.com/77336594/226332268-93b8fa85-ba5b-4fcc-938b-ca8d642b8521.png)

7. Navigate to **API permissions** in the left pane.

8. Click on **Add a permission**, then **Microsoft Graph** in the right pane.

![image](https://user-images.githubusercontent.com/6913320/76350226-c2961200-62e1-11ea-9080-19a9b75c2aee.png)

9. Click on **Delegated permissions**, and scroll down to select the following permissions:

- `ChannelMessage.Send`
- `ChatMessage.Send`
- `email`
- `offline_access`
- `openid`
- `profile`
- `User.Read`

10. Click on **Add permissions** to submit the form.

11. Next, add application permissions via **Add a permission > Microsoft Graph > Application permissions**.

12. Select the following permissions:

- `Chat.Create`
- `Group.Create`
- `Group.Read.All`
- `Group.ReadWrite.All`
- `Team.ReadBasic.All`
- `TeamMember.ReadWrite.All`
- `TeamMember.ReadWriteNonOwnerRole.All`

13. Click on **Add permissions** to submit the form.

14. Click on **Grant admin consent for...** to grant the permissions for the application.

15. Navigate to **Authentication** in the left pane.

16. Scroll down below to **Advanced settings** and change **Allow public client flows value** to `Yes`.

![image](https://github.com/Brightscout/msteams-load-test-scripts/assets/77336594/5a759c30-9f73-4570-b201-6a183a567691)
