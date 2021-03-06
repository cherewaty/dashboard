// Copyright 2017 The Kubernetes Dashboard Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

import {ReplicaSetListController} from 'replicaset/list/controller';
import replicaSetModule from 'replicaset/module';

describe('Replica Set list controller', () => {
  /** @type {!ReplicaSetListController} */
  let ctrl;

  beforeEach(() => {
    angular.mock.module(replicaSetModule.name);

    angular.mock.inject(($controller) => {
      ctrl = $controller(ReplicaSetListController, {replicaSetList: {replicaSets: []}});
    });
  });

  it('should initialize replica set controller', angular.mock.inject(($controller) => {
    let ctrls = {};
    /** @type {!ReplicaSetListController} */
    let ctrl = $controller(ReplicaSetListController, {replicaSetList: {replicaSets: ctrls}});

    expect(ctrl.replicaSetList.replicaSets).toBe(ctrls);
  }));

  it('should show zero state', () => {
    expect(ctrl.shouldShowZeroState()).toBeTruthy();
  });

  it('should hide zero state', () => {
    // given
    ctrl.replicaSetList = {replicaSets: ['mock']};

    // then
    expect(ctrl.shouldShowZeroState()).toBeFalsy();
  });
});
