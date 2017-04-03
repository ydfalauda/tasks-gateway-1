angular.module('TasksApp')
.service('ToastDisplay',['$mdToast', function ($mdToast) {
  return {
    showToast: function (error) {
      $mdToast.show(
        $mdToast.simple()
          .textContent(error)
          .position('top')
          .hideDelay(4000)
      )
    }
  }
}])


/**
 * Login/Signup Controller for the TasksApp
 * @param UserService
 * @param $location
 * @constructor
 */
angular.module('TasksApp')
.controller('LoginCtrl', ['UserService','$location', 'ToastDisplay',  LoginController])

function LoginController (UserService, $location, ToastDisplay) {
  var self = this;

  //Extending the controller
  angular.extend(self, ToastDisplay);

  self.submitLogin = submitLogin;
  self.submitSignup = submitSignup;
  self.credentials = {
    username: "",
    password: ""
  };

  //Submit login function
  function submitLogin () {
    UserService.login(self.credentials)
      .then(function (token) {
        console.log(token);
        // alert('got token bitch');
        $location.path('/tasks');
      },
      self.showToast);
  }

  //Signup function
  function submitSignup () {
    UserService.signup(self.credentials)
      .then(function (token) {
        console.log(token);
        // alert('got token bitch');
        $location.path('/tasks');
      },
      self.showToast);
  }

}

/**
 * Tasks Controller for the TasksApp
 * @param TasksService
 * @param $location
 * @param $mdDialog
 * @param ToastDisplay
 * @constructor
 */
angular.module('TasksApp')
.controller('TasksCtrl', ['TasksService', '$location', '$mdDialog' ,'ToastDisplay', TasksController])

function TasksController (TasksService, $location, $mdDialog, ToastDisplay) {
  var self = this;

  //Extending
  angular.extend(self, ToastDisplay);

  self.items    = [];
  self.selected = undefined;
  self.backup   = null;
  self.loaded   = false;

  self.refresh = function () {
    TasksService.all()
      .then( function (tasks) {
        self.items = [].concat(tasks);
        self.loaded = true;
      }, self.showToast);
  }

  self.refresh();

  self.isTaskSelected = function (task) {
    if (self.selected === null || self.selected === undefined) return false;

    return (self.selected.id && self.selected.id === task.id);
  }

  self.addNew = function (event) {
    var confirm = $mdDialog.prompt()
          .title('What is the new task?')
          .placeholder('task')
          .ariaLabel('Task name')
          .targetEvent(event)
          .ok('Add')
          .cancel('Cancel');
    $mdDialog.show(confirm).then(function(result) {

      if (result && result.trim()) {
        TasksService.add(result)
          .then(function (added) {
            self.items.unshift(added);
          },
          self.showToast)
      }
    });
  }

  self.checkTask = function (task, event) {
    event.stopImmediatePropagation();

    //Check event, reverting the done status
    task.done = !task.done;

    TasksService.update(task)
      .then(function (added) {
        self.showToast(added);
      },
      function (error){
        //Reverting state
        task.done = !task.done;
        //Showing error
        self.showToast(error);
      })
  }


  self.updateTask = function (task) {
    TasksService.update(task)
      .then(function (added) {
      },
      self.showToast)
  }


  self.edit = function (task, event) {
    if (!task) return;

    event.stopImmediatePropagation();

    //Saving the name of the task
    self.backup = task.name;

    var confirm = $mdDialog.prompt()
          .title('Edit the task?')
          .placeholder('task')
          .textContent('Current: '+task.name)
          .ariaLabel('Task name')
          .targetEvent(event)
          .ok('Add')
          .cancel('Cancel');
    $mdDialog.show(confirm).then(function(result) {
      if (!result || !result.trim())
        return;
      task.name = result;
      TasksService.update(task)
        .then(function (added) {
        },
        function (error){
          //Reverting the change
          task.name = self.backup;
          //Displaying the error
          self.showToast(error);
        });
    });
  }
}

/**
 * Logout Controller for the TasksApp
 * @param UserService
 * @param $location
 * @param ToastDisplay
 * @constructor
 */
angular.module('TasksApp')
.controller('LogoutCtrl', ['UserService', '$location', 'ToastDisplay', LogoutController])

function LogoutController (UserService, $location, ToastDisplay) {
  var self = this;

  angular.extend(self, ToastDisplay);

  self.logoutUser = logoutUser;

  function logoutUser () {
    if (UserService.getToken()) {
      UserService.logout()
        .then(function (ok) {
          $location.path('/');
        },
        self.showToast);
    }
  }
}
