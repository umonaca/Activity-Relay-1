/*
Models provide struct for config, state and type definition for ActivityPub, Nodeinfo, Webfinger.

Redis schema of RelayState

RelayState is permanent state information backed redis key value store.

RelayState.relayOptions

RelayState.relayOptions are associated with below.

	hash relay:config
		field block_service      - value 0(false) or 1(true)
		field manually_accept    - value 0(false) or 1(true)
		field create_as_announce - value 0(false) or 1(true)

RelayState.limitedDomains

RelayState.limitedDomains are associated with below.

	hash relay:config:limitedDomain
		field <DomainHostName> - value 1(true)

RelayState.blockedDomains

RelayState.blockedDomains are associated with below.

	hash relay:config:blockedDomain
		field <DomainHostName> - value 1(true)

RelayState.subscriptions

RelayState.subscriptions are associated with below.

	hash relay:subscription:<DomainHostName>
		field inbox_url   - value <Inbox URL>
		field activity_id - value <Follow Activity's ID>
		field actor_id    - value <Follow Actor's ID>
		field available   - value true(Subscribe) or false(Pending)

RelayState.StateRefreshListener

RelayState.StateRefreshListener is state change notification for cluster.

	channel relay_refresh
		message - value nil
*/
package models
