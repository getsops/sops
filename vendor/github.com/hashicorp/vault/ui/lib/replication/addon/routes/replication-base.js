import { alias } from '@ember/object/computed';
import { inject as service } from '@ember/service';
import { hash, resolve } from 'rsvp';
import Route from '@ember/routing/route';
import UnloadModelRouteMixin from 'vault/mixins/unload-model-route';

export default Route.extend(UnloadModelRouteMixin, {
  store: service(),
  version: service(),
  rm: service('replication-mode'),
  modelPath: 'model.config',
  fetchMounts() {
    return hash({
      mounts: this.store.findAll('secret-engine'),
      auth: this.store.findAll('auth-method'),
    }).then(({ mounts, auth }) => {
      return resolve(mounts.toArray().concat(auth.toArray()));
    });
  },

  replicationMode: alias('rm.mode'),
});
