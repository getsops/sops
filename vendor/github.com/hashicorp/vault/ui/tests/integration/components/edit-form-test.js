import { later, run } from '@ember/runloop';
import { resolve } from 'rsvp';
import EmberObject from '@ember/object';
import Service from '@ember/service';
import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, settled } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import sinon from 'sinon';
import { create } from 'ember-cli-page-object';
import editForm from 'vault/tests/pages/components/edit-form';

const component = create(editForm);

const flash = Service.extend({
  success: sinon.stub(),
});

const createModel = (canDelete = true) => {
  return EmberObject.create({
    fields: [{ name: 'one', type: 'string' }, { name: 'two', type: 'boolean' }],
    canDelete,
    destroyRecord() {
      return resolve();
    },
    save() {
      return resolve();
    },
    rollbackAttributes() {},
  });
};

module('Integration | Component | edit form', function(hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function() {
    run(() => {
      this.owner.unregister('service:flash-messages');
      this.owner.register('service:flash-messages', flash);
    });
  });

  test('it renders', async function(assert) {
    let model = createModel();
    this.set('model', model);
    await render(hbs`{{edit-form model=model}}`);

    assert.ok(component.fields.length, 2);
  });

  test('it calls flash message fns on save', async function(assert) {
    let model = createModel();
    let onSave = () => {
      return resolve();
    };
    this.set('model', model);
    this.set('onSave', onSave);
    let saveSpy = sinon.spy(this, 'onSave');

    await render(hbs`{{edit-form model=model onSave=onSave}}`);

    component.submit();
    later(() => run.cancelTimers(), 50);
    return settled().then(() => {
      assert.ok(saveSpy.calledOnce, 'calls passed onSave');
      assert.equal(saveSpy.getCall(0).args[0].saveType, 'save');
      assert.deepEqual(saveSpy.getCall(0).args[0].model, model, 'passes model to onSave');
      let flash = this.owner.lookup('service:flash-messages');
      assert.equal(flash.success.callCount, 1, 'calls flash message success');
    });
  });
});
