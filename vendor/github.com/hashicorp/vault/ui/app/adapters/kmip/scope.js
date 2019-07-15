import BaseAdapter from './base';

export default BaseAdapter.extend({
  createRecord(store, type, snapshot) {
    let name = snapshot.attr('name');
    return this.ajax(this._url(type.modelName, { backend: snapshot.record.backend }, name), 'POST').then(
      () => {
        return {
          id: name,
          name,
        };
      }
    );
  },

  deleteRecord(store, type, snapshot) {
    return this.ajax(this._url(type.modelName, { backend: snapshot.record.backend }, snapshot.id), 'DELETE');
  },
});
