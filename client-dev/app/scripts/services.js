//User Service
angular.module('TasksApp')
.service('UserService', ['$http', 'apiUrl', 'localStorageService', UserService])

function UserService ($http, apiUrl, localStorageService) {
  var self = this;

  var userUrl = apiUrl+'/users';

  function setHeader (token) {
    $http.defaults.headers.common.Authorization = token;
  }
  function getHeader () {
    return $http.defaults.headers.common.Authorization;
  }

  function loginUser (credentials) {
    return $http.post( userUrl+'/login', credentials).then(
      function (response) {
        var token = response.data;

        //Save token in the localstorage
        localStorageService.set('at', token.token);
        //Set Authorization header
        setHeader(token.token);

        return token;
      },
      throwError
    );
  }

  function signupUser (credentials) {
    return $http.post( userUrl+'/signup', credentials).then(
      function (response) {
        var token = response.data;

        //Save token in the localstorage
        localStorageService.set('at', token.token);
        //Set Authorization header
        setHeader(token.token);

        return token;
      },
      throwError
    );
  }

  function getUserToken () {
    var token = localStorageService.get('at');
    if (token !== null && token !== undefined && (!getHeader()|| getHeader() !== token )) {
      setHeader(token);
    }
    return token !== null && token !== undefined;
  }

  function logout() {
    return $http.post( userUrl+'/logout').then(
      function () {
        localStorageService.remove('at');
        setHeader(undefined);
        return true;
      },
      function () {
        localStorageService.remove('at');
        setHeader(undefined);
        return true;
      }
    );
  }

  return {
    login: loginUser,
    signup: signupUser,
    getToken: getUserToken,
    logout: logout
  }
}

angular.module('TasksApp')
.service('TasksService', ['$http', 'apiUrl','UserService', TasksService])

/**
 * Task Service for Task App.
 * @param $http
 * @param apiUrl
 * @param UserService
 * @constructor
 */

function TasksService ($http, apiUrl, UserService) {
  var self = this;
  var tasksUrl = apiUrl+'/tasks';

  function getAllTasks () {
    return $http.get(tasksUrl)
      .then(function (response) {
        return response.data;
      }, throwError);
  }

  function addNewTask (name) {
    return $http.post(tasksUrl, {name: name})
      .then(function (response) {
        return response.data;
      }, throwError);
  }

  function updateTask (task) {
    var data = {};
    if (task.id) data.id = task.id;
    if (task.name) data.name = task.name;
    if (task.done !== null && task.done !== undefined) data.done = task.done;

    return $http.put(tasksUrl, data)
      .then(function (response) {
        return response.data;
      }, throwError);
  }

  return {
    all: getAllTasks,
    add: addNewTask,
    update: updateTask
  };
}

/**
 * throwError Helper for Services.
 * @param httpError
 * @constructor
 */
function throwError (httpError) {
  throw httpError.status + " : " + httpError.data;
}
