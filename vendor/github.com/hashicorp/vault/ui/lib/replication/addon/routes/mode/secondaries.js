import { setProperties } from '@ember/object';
import { hash } from 'rsvp';
import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';

export default Route.extend({
  store: service(),
  model() {
    const replicationMode = this.paramsFor('mode').replication_mode;

    return hash({
      cluster: this.modelFor('mode'),
      canAddSecondary: this.store
        .findRecord('capabilities', `sys/replication/${replicationMode}/primary/secondary-token`)
        .then(c => c.get('canUpdate')),
      canRevokeSecondary: this.store
        .findRecord('capabilities', `sys/replication/${replicationMode}/primary/revoke-secondary`)
        .then(c => c.get('canUpdate')),
    }).then(({ cluster, canAddSecondary, canRevokeSecondary }) => {
      setProperties(cluster, {
        canRevokeSecondary,
        canAddSecondary,
      });
      return cluster;
    });
  },
  afterModel(model) {
    const replicationMode = this.paramsFor('mode').replication_mode;
    if (
      !model.get(`${replicationMode}.isPrimary`) ||
      model.get(`${replicationMode}.replicationDisabled`) ||
      model.get(`${replicationMode}.replicationUnsupported`)
    ) {
      return this.transitionTo('mode', replicationMode);
    }
  },
});
