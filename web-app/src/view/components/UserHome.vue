<!--suppress ALL -->
<template>
    <div>
        <div class="nav-wrp">
            <navigation-bar />
        </div>
        <div class="columns">
            <div class="column">
                <user-badge :user="userData" />
                <div class="box">
                    <p class="title is-6">Upload User Foto</p>
                    <div class="file-upload-wrp">
                        <b-upload v-model="userPhoto.file">
                            <a class="button is-light">
                                <b-icon icon="upload"></b-icon>
                            </a>
                        </b-upload>
                        <span class="file-name">
                            {{ userPhoto.file ? userPhoto.file.name : 'Pilih file' }}
                        </span>
                    </div>
                    <a @click="doSave(true)" v-if="userPhoto.file" class="bt-upload button is-white">Upload</a>
                </div>
                <div class="box">
                    <p class="title is-6">User Foto</p>
                    <template v-if="userPhotos.length > 0">
                        <div class="photo-view">
                            <p v-for="(userPhoto, index) in userPhotos" :key="'user-photo-'+index" class="image is-cover">
                                <common-image :src="userPhoto.path" />
                            </p>
                        </div>
                        <a @click="() => hasMoreUserPhoto ? doFind(userPhotos.length, 'USER_PHOTO') : null" class="button is-white bt-load">
                            {{hasMoreUserPhoto ? 'Cari' : 'Tak ada lagi konten dimuat'}}
                        </a>
                    </template>
                    <empty-states v-else :isVertical="true" />
                </div>
                <div class="box has-text-centered">
                    <a class="button is-white" href="/admin" data-turbolinks-action="replace">Lihat Admin</a>
                </div>
            </div>
            <div class="column is-two-thirds">
                <div class="box">
                    <p class="title is-6">{{(selectedUser || {username:''}).username.length > 0 ? `Selected user is ${(selectedUser || {username:''}).username}` : 'Example autocomplete buefy'}}</p>
                    <b-field label="Find a user">
                        <b-autocomplete v-model="username" :data="userDatas" placeholder="e.g. Will smith" field="username" :loading="isFetching" @keyup.native="onTypingUser" @select="option => selectedUser = option">
                            <template slot-scope="props">
                                <div class="media">
                                    <div class="media-left">
                                        <common-image width="32" :src="props.option.image_profile" :isProfile="true" />
                                    </div>
                                    <div class="media-content">
                                        {{ props.option.username }}
                                        <br>
                                        <small>
                                            Name {{props.option.name}}
                                            Created at {{ props.option.created_at }},
                                        </small>
                                    </div>
                                </div>
                            </template>
                            <template slot="empty">
                                <empty-states />
                            </template>
                            <template slot="header">
                                <a @click="doFind(userDatas.length, 'USER')">
                                    <span> Load more </span>
                                </a>
                            </template>
                        </b-autocomplete>
                    </b-field>
                </div>

                <common-editor v-model="todo.content" ref="commonEditor" @update="({html})=> todo.content = html">
                    <template slot="customMenubar">
                        <b-upload @input="onUploadToEditor" v-model="todo.todoFileImage">
                            <button class="button is-white">
                                <i class="mdi mdi-image-area"></i>
                            </button>
                        </b-upload>
                    </template>
                    <b-field slot="slot-before" label="Name" :type="todo.nameFeedback().type" :message="todo.nameFeedback().error">
                        <b-input v-model="todo.name" size="is-small" placeholder="e.g Kill Your self" maxlength="50"></b-input>
                    </b-field>
                    <template slot="slot-after">
                        <b-field label="Select a date">
                            <b-datepicker :min-date="todo.minDateDue" v-model="todo.due" size="is-small" placeholder="Click to select..." icon="calendar-today">
                            </b-datepicker>
                        </b-field>
                        <b-field>
                            <p class="control">
                                <a @click="doSave(false)" :class="`button is-${todo.valid() ? 'info' : 'danger'} has-icon`">
                                    <span class="icon">
                                        <i class="mdi mdi-content-save"></i>
                                    </span>
                                    <span>Save</span>
                                </a>
                            </p>
                        </b-field>
                    </template>
                </common-editor>
                <div class="todo-item box" v-for="(todo, index) in todos" :key="'todo-item-'+index">
                    <div class="content">
                        <div>
                            <div class="level">
                                <div class="level-left">
                                    <p class="title is-5">
                                        <a :href="`/todo/${todo.name}`" data-turbolinks-action="replace">{{todo.name}}</a>
                                    </p>
                                </div>
                                <div class="level-right">
                                    <a @click="doEdit(todo)" class="button is-small">
                                        <span>
                                            <i class="mdi mdi-circle-edit-outline"></i>
                                        </span>
                                    </a>
                                </div>
                            </div>
                            <div class="content-todo">
                                <p>{{todo.subcontent}}</p>
                            </div>
                            <div class="field">
                                <p class="help is-grey">{{todo.user.name}}</p>
                                <span>•</span>
                                <p class="help is-grey">{{todo.created_at}}</p>
                            </div>
                        </div>
                    </div>
                </div>
                <a @click="() => hasMoreTodo ? doFind(todos.length) : null" class="button is-white bt-load">
                    {{hasMoreTodo ? 'Cari' : 'Tak ada lagi konten dimuat'}}
                </a>
            </div>
        </div>
    </div>
</template>

<script lang="ts">
import { Vue, Component, Prop, Inject } from "annotation";
import { Deserialize } from "cerialize";
import environment from "environment";
import { isEmpty, isNil, debounce, uniqBy } from "lodash";
import Constant from "../config/Constant";

import TodoService from "../service/TodoService";
import UserService from "../service/UserService";
import UserPhotoService from "../service/UserPhotoService";

import User from "../models/User";
import Todo from "../models/Todo";
import UserPhoto from "../models/UserPhoto";

@Component
export default class UserHome extends Vue {
  @Prop({
    default: ""
  })
  private user: string;

  @Inject
  private todoService: TodoService;

  @Inject
  private userPhotoService: UserPhotoService;

  @Inject
  private userService: UserService;

  private hasMoreTodo: boolean = true;

  private isFetching: boolean = false;

  private username: string = "";

  private hasMoreUserPhoto: boolean = true;

  private todo: Todo = new Todo();

  private todos: Array<Todo> = [];

  private userPhoto: UserPhoto = new UserPhoto();

  private userPhotos: Array<UserPhoto> = [];

  private userData: User = new User();

  private userDatas: Array<User> = [];

  private selectedUser: User = new User();

  private created() {
    try {
      const userJson = JSON.parse(this.user);

      this.userData = Deserialize(userJson, User);
      this.userPhoto.user_id = this.todo.user_id = this.userData.id;
    } catch (err) {
      console.log(err);
    }
  }

  private async mounted() {
    await this.doFind();
    await this.doFind(0, "USER_PHOTO");
  }

  private async doFind(offset: number = 0, type: string = "TODOS") {
    this.isFetching = true;

    const user_id = this.userData.id;
    const param = {
      offset,
      limit: environment["LIMIT"],
      user_id
    };

    switch (type) {
      case "USER_PHOTO":
        const userPhotos = await this.userPhotoService.find(param);
        if (offset === 0) {
          this.userPhotos = [];
        }

        this.userPhotos = this.userPhotos.concat(
          Array.isArray(userPhotos) ? userPhotos : []
        );
        this.hasMoreUserPhoto =
          userPhotos.length % environment["LIMIT"] === 0 &&
          userPhotos.length !== 0;

        break;
      case "USER":
        const paramUser = {
          offset: 0,
          limit: environment["LIMIT"],
          username: this.username
        };
        const users = await this.userService.find(paramUser);
        if (offset === 0) {
          this.userDatas = [];
        }

        this.userDatas = uniqBy(
          this.userDatas.concat(Array.isArray(users) ? users : []),
          "id"
        );

        break;
      default:
        const todos = await this.todoService.find(param);
        if (offset === 0) {
          this.todos = [];
        }

        this.todos = this.todos.concat(Array.isArray(todos) ? todos : []);
        this.hasMoreTodo =
          todos.length % environment["LIMIT"] === 0 && todos.length !== 0;

        break;
    }

    this.isFetching = false;
  }

  private onUploadToEditor(files: any) {
    const userPhoto = new UserPhoto();
    userPhoto.user_id = this.userData.id;
    userPhoto.file = files;

    this.doSave(true, userPhoto);
  }

  private resetEntity(onlyResetUserPhoto: boolean = false) {
    this.userPhoto = new UserPhoto();

    if (!onlyResetUserPhoto) {
      this.todo = new Todo();
      (this as any).$refs.commonEditor.resetEditor();
    }

    this.userPhoto.user_id = this.todo.user_id = this.userData.id;
  }

  private async doSave(isUpload: boolean = false, userPhoto: UserPhoto = null) {
    this.todoService.returnWithStatus = true;
    this.userPhotoService.returnWithStatus = true;
    let response = null;

    if (!isNil(userPhoto)) {
      response = await this.userPhotoService.save(userPhoto);
      (this as any).$refs.commonEditor.setContent(
        this.todo.content +
          `<p class="has-text-centered"><img src="${response.data.path}" /></p>`
      );
    } else if (this.userPhoto.file && isUpload) {
      response = await this.userPhotoService.save(this.userPhoto);
    } else if (!isUpload && this.todo.valid()) {
      if (this.todo.id > 0) {
        response = await this.todoService.update(
          this.todo,
          this.todo.id.toString()
        );
      } else {
        response = await this.todoService.save(this.todo);
      }
    }

    const { SAVE_SUCCESS, UPDATE_SUCCESS } = Constant.STATUS.API;
    if (
      response.status === SAVE_SUCCESS ||
      response.status === UPDATE_SUCCESS
    ) {
      this.doFind(0, "USER_PHOTO");
      this.resetEntity(!isNil(userPhoto));
    }

    this.todoService.returnWithStatus = false;
    this.userPhotoService.returnWithStatus = false;
  }

  private async doEdit(todo: Todo) {
    const response = await this.todoService.find("/" + todo.id);
    if (response) {
      this.todo = response;
    } else {
      this.todo = todo;
    }

    (this as any).$refs.commonEditor.setContent(this.todo.content);
  }

  private onTypingUser: any = debounce(function() {
    if (!this.username.length) {
      this.userDatas = [];

      return;
    }

    this.doFind(0, "USER").then();
  });
}
</script>
