import com.onresolve.scriptrunner.runner.customisers.WithPlugin
import com.onresolve.scriptrunner.runner.customisers.PluginModule
import com.atlassian.rm.teams.api.team.GeneralTeamService
import com.atlassian.jira.component.ComponentAccessor

@WithPlugin("com.atlassian.teams")
@PluginModule
GeneralTeamService teamService

def userManager = ComponentAccessor.userManager
def teams = teamService.getAllShareableTeams()

def csv = new StringBuilder()
csv << "teamId,teamName,memberCount,userKey,username,displayName,active,resourceId,personId\n"

teams.each { teamId, team ->
    def teamName = team.getDescription()?.getTitle() ?: ""
    def resources = team.getResources() ?: []

    if (!resources) {
        csv << [
            teamId,
            teamName,
            0,
            "",
            "",
            "",
            "",
            "",
            ""
        ].collect { "\"${(it ?: '').toString().replace('"', '""')}\"" }.join(",") << "\n"
        return
    }

    resources.each { resource ->
        def person = resource?.getPerson()
        def personDescription = person?.getDescription()

        def jiraUserKey = null
        try {
            jiraUserKey = personDescription?.getJiraUser()?.orNull()
        } catch (ignored) {
            jiraUserKey = personDescription?.getJiraUser()?.toString()
        }

        def jiraUser = jiraUserKey ? userManager.getUserByKey(jiraUserKey as String) : null

        csv << [
            teamId,
            teamName,
            resources.size(),
            jiraUserKey,
            jiraUser?.name,
            jiraUser?.displayName,
            jiraUser?.active,
            resource?.getId(),
            person?.getId()
        ].collect { "\"${(it ?: '').toString().replace('"', '""')}\"" }.join(",") << "\n"
    }
}

return csv.toString()