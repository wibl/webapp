const EditMode = {
  CREATE: Symbol(), CHANGE: Symbol(), DELETE: Symbol()
}

const api = {
  async invoke(method, params) {
    const response = await axios.post("/api", {"jsonrpc": "2.0", "method": method, "params": [params || {}], "id": 1});
    if (response.data.error) {
      throw new Error(response.data.error);
    }
    return response.data.result;
  },
  async connectToMq(url, user, pass) {
    return this.invoke("MqService.Connect", {"URL": url, "User": user, "Pass": pass});
  },
  async getAllGroups() {
    return this.invoke("GS.GetAllGroups");
  },
  async createGroup(title) {
    return this.invoke("GS.CreateGroup", {"Title": title});
  },
  async deleteGroup(groupId) {
    return this.invoke("GS.DeleteGroup", {"ID": groupId});
  },
  async getAllTemplates(groupId) {
    return this.invoke("TS.GetAllTemplates", {"GroupID": groupId});
  },
  async createTemplate(groupId, title, queue, body) {
    return this.invoke("TS.CreateTemplate", {"GroupID": groupId, "Title": title, "Queue": queue, "Body": body});
  },
  async sendToMq(queue, message) {
    return this.invoke("MqService.Send", {"Queue": queue, "Message": message});
  }
}

const store = new Vuex.Store({
  state: {
    groups: [],
    templates: [],
    selectedGroupId: 0,
    selectedTemplateId: 0
  },
  mutations: {
    addGroup (state, group) {
      state.groups.push(group);
    },
    deleteGroup (state, group) {
      state.groups.splice(state.groups.indexOf(group), 1);
      state.templates.filter(template => template.GroupID === group.ID).forEach(template => {
        state.templates.splice(state.templates.indexOf(template), 1);
      });
    },
    addTemplate (state, template) {
      state.templates.push(template);
    },
    setSelectedGroupId (state, groupId) {
      state.selectedGroupId = groupId;
    },
    setSelectedTemplateId (state, templateId) {
      state.selectedTemplateId = templateId;
    },
  },
  actions: {
    async signIn (context, payload) {
      await api.connectToMq(payload.url, payload.user, payload.pass);

      const resAllGroups = await api.getAllGroups();
      if (resAllGroups.Groups) {
        for (i = 0; i < resAllGroups.Groups.length; i++) {
          let group = resAllGroups.Groups[i];
          context.commit("addGroup", group);
          
          const resAllTemplates = await api.getAllTemplates(group.ID);
          if (resAllTemplates.Templates) {
            for (j = 0; j < resAllTemplates.Templates.length; j++) {
              let template = resAllTemplates.Templates[j]
              context.commit("addTemplate", template);
              
              if (i === 0 && j === 0) {
                context.commit("setSelectedGroupId", group.ID);
                context.commit("setSelectedTemplateId", template.ID);
              }
            }
          }
        }
      }
    },
    async save (context, payload) {
      let groupId = payload.groupId;

      switch (payload.groupEditMode) {
        case EditMode.CREATE:
          const result = await api.createGroup(payload.groupTitle);

          context.commit("addGroup", result.Group);
          groupId = result.Group.ID;
          break;
        case EditMode.DELETE:
          await api.deleteGroup(groupId);

          context.commit("deleteGroup", groupId);
          break;
        case EditMode.CHANGE:
          //TODO: Implement
          break;
      }

      switch (payload.templateEditMode) {
        case EditMode.CREATE:
            const result = await api.createTemplate(groupId, payload.templateTitle, payload.queue, payload.body);

            context.commit("addTemplate", result.Template);
        break;
        case EditMode.DELETE:
          //TODO: Implement
          break;
        case EditMode.CHANGE:
          //TODO: Implement
          break;
      }

    },
    async send (context, payload) {
      await api.sendToMq(payload.queue, payload.message);
    }
  },
  getters: {
    getTemplatesByGroupId: (state) => (groupId) => {
      let templates = state.templates.filter(template => template.GroupID === groupId);
      console.log("getTemplatesByGroupId groupId", groupId, "templates", templates);
      console.log("getTemplatesByGroupId state.templates", state.templates);
      return templates;
    },
    getTemplateById: (state) => (templateId) => {
      return state.templates.find(template => template.ID === templateId);
    },
    getGroupById: (state) => (groupId) => {
      let groups = state.groups.find(group => group.ID === groupId);
      return groups;
    }
  }
})

const LoginPage = {
  template: `
    <form @submit.prevent="onSubmit" class="pure-form pure-form-aligned">
      <div v-if="errors.length" class="alert alert-danger">
        <ul v-for="error in errors" class="list-unstyled">
          <li>{{ error }}</li>
        </ul>
      </div>
    	<fieldset>
        <div class="pure-control-group">
          <label for="url">URL:</label>
          <input id="url" v-model="url">
        </div>
        <div class="pure-control-group">
          <label for="user">User(opt):</label>
          <input id="user" v-model="user">
        </div>
        <div class="pure-control-group">
          <label for="pass">Pass(opt):</label>
          <input id="pass" type="password" v-model="pass">
        </div>
        <div class="pure-controls">
          <button type="submit" :disabled="!isSubmitButtonEnabled" class="pure-button pure-button-primary">Sign in</button>
        </div>
      </fieldset>
		</form>
	`,
  methods: {
    onSubmit: function() {
      this.isSubmitButtonEnabled = false;
      this.errors = [];

      this.$store.dispatch("signIn", {url: this.url, user: this.user, pass: this.pass}).then(() => {
        this.$router.push("/main");
      }).catch((err) => {
        this.errors.push(err)
        this.isSubmitButtonEnabled = true;
      });
      
    }
  },
  data: function() {
    return {
      url: "tcp://localhost:61613",
      user: "",
      pass: "",
      isSubmitButtonEnabled: true,
      errors: []
    }
  }
}

const MainPage = {
  template: `
  	<form @submit.prevent="onSubmit" class="pure-form pure-form-aligned">
    	<fieldset>
        <div class="pure-control-group">
          <label for="group">Group:</label>
          <select id="group" :value="selectedGroupId" @input="setSelectedGroupId">
            <option v-for="group in groupList" v-bind:value="group.ID">
              {{ group.Title }}
            </option>
          </select>
        </div>
        <div class="pure-control-group">
          <label for="template">Template:</label>
          <select id="template" :value="selectedTemplateId" @input="setSelectedTemplateId">
            <option v-for="template in templateList" v-bind:value="template.ID">
              {{ template.Title }}
            </option>
          </select>
          <button @click.prevent="onEdit" class="pure-button">Edit</button>
        </div>
        <div class="pure-control-group">
        	<label for="queue">Queue:</label>
          <input id="queue" type="text" :value="selectedTemplateQueue" readonly>
        </div>
        <div class="pure-control-group">
          <label for="message">Message:</label>
          <textarea id="message" :value="selectedTemplateBody" style="vertical-align: top;" />
        </div>
        <div class="pure-controls">
          <button type="submit" class="pure-button pure-button-primary">Send</button>
        </div>
      </fieldset>
    </form>
	`,
  methods: {
    onSubmit() {
      this.$store.dispatch("send", {queue: this.selectedTemplateQueue, message: this.selectedTemplateBody}).then(() => {
        alert("Message sent")
      })
    },
    onEdit() {
      if (this.selectedGroupId === 0) {
        this.$router.push({name: "create_group"});
      } else if (this.selectedTemplateId === 0) {
        this.$router.push({
          name: "create_template",
          params: {
            groupId: this.selectedGroupId
          }
        });
      } else {
        this.$router.push({
          name: "edit_template",
          params: {
            groupId: this.selectedGroupId,
            templateId: this.selectedTemplateId
          }
        });
      }
    },
    setSelectedGroupId(event) {
      this.$store.commit("setSelectedGroupId", event.target.value);
    },
    setSelectedTemplateId(event) {
      this.$store.commit("setSelectedTemplateId", event.target.value);
    },
  },
  data() {
    return {
      // firstGroup: this.$store.state.groups[0],
      // firstTemplate: this.firstGroup ? this.$store.getters.getTemplatesByGroupId(this.firstGroup.ID)[0] : "",
      // selectedGroupId: this.$store.state.groups[0] ? this.$store.state.groups[0].ID : 0,
      // selectedTemplateId: selectedGroupId ? ,
    }
  },
  computed: {
    selectedGroupId() {
      let groupId = this.$store.state.selectedGroupId;
      console.log("selectedGroupId", groupId)
      return groupId;
    },
    groupList() {
      let groups = this.$store.state.groups;
      console.log("groupList", groups);
      return groups;
    },
    templateList() {
      let templates = this.$store.getters.getTemplatesByGroupId(this.selectedGroupId);
      console.log("templateList groupId", this.selectedGroupId, "templates", templates);
      return templates;
    },
    selectedTemplateId() {
      let templateId = this.$store.state.selectedTemplateId;
      console.log("selectedTemplateId", templateId);
      return templateId;
    },
    selectedTemplate() {
      return this.$store.getters.getTemplateById(this.selectedTemplateId);
    },
    selectedTemplateQueue() {
      return this.selectedTemplate ? this.selectedTemplate.Queue : "";
    },
    selectedTemplateBody() {
      return this.selectedTemplate ? this.selectedTemplate.Body : "";
    }
  }
}

const EditPage = {
  props: ["groupId", "templateId"],
  template: `
    <form class="pure-form pure-form-aligned">
      groupEditMode=<div v-if="groupEditMode === EditMode.CREATE">CREATE</div><div v-if="groupEditMode === EditMode.CHANGE">CHANGE</div><div v-if="groupEditMode === EditMode.DELETE">DELETE</div>
      templateEditMode=<div v-if="templateEditMode === EditMode.CREATE">CREATE</div><div v-if="templateEditMode === EditMode.CHANGE">CHANGE</div><div v-if="templateEditMode === EditMode.DELETE">DELETE</div>
      <div v-if="errors.length" class="alert alert-danger">
        <ul v-for="error in errors" class="list-unstyled">
          <li>{{ error }}</li>
        </ul>
      </div>
      <fieldset>
        <div class="pure-control-group">
          <label for="groupTitle">Group title:</label>
          <input id="groupTitle" v-model="groupTitle">
          <button @click.prevent="onGroupAdd" class="pure-button">Add</button>
          <button @click.prevent="onGroupDelete" class="pure-button">Delete</button>
				</div>
        <div class="pure-control-group">
          <label for="templateTitle">Template title:</label>
          <input id="templateTitle" v-model="templateTitle">
          <button @click.prevent="onTemplateAdd" class="pure-button">Add</button>
          <button @click.prevent="onTemplateDelete" class="pure-button">Delete</button>
        </div>
        <div class="pure-control-group">
          <label for="queue">Queue:</label>
          <input id="queue" v-model="queue">
        </div>
        <div class="pure-control-group">
          <label for="body">Message:</label>
          <input id="body" v-model="body">
        </div>
        <div class="pure-controls">
        	<button @click.prevent="onSave" class="pure-button pure-button-primary">Save</button>
          <button @click.prevent="onCancel" class="pure-button">Cancel</button>
        </div>
      </fieldset>
      <div v-if="isSaving" class="overlay-loader"></div>
    </form>
	`,
  methods: {
    onSave: function() {
      this.isSaving = true;

      this.$store.dispatch("save", {
        groupEditMode: this.groupEditMode,
        templateEditMode: this.templateEditMode,
        groupId: this.groupId,
        groupTitle: this.groupTitle,
        templateId: this.templateId,
        templateTitle: this.templateTitle,
        queue: this.queue,
        body: this.body
      }).then(() => {
        this.$router.push("/main");
      }).catch((err) => {
        this.errors.push(err);
        this.isSaving = false;
      });
    },
    onCancel: function() {
      this.$router.push("/main");
    },
    onGroupAdd: function() {
      this.groupEditMode = EditMode.CREATE;
    },
    onGroupDelete: function() {
      this.groupEditMode = EditMode.DELETE;
    },
    onTemplateAdd: function() {
      this.templateEditMode = EditMode.CREATE;
    },
    onTemplateDelete: function() {
      this.templateEditMode = EditMode.DELETE;
    },
  },
  data: function() {
    return {
      groupEditMode: EditMode.CREATE,
      templateEditMode: EditMode.CREATE,
      groupTitle: "",
      templateTitle: "",
      queue: "",
      body: "",
      EditMode: EditMode,
      isSaving: false,
      errors: []
    }
  },
  created: function() {
    if (this.groupId) {
      this.groupEditMode = EditMode.CHANGE

      const editableGroup = this.$store.getters.getGroupById(this.groupId);
      this.groupTitle = editableGroup.Title
    }
    if (this.templateId) {
      this.templateEditMode = EditMode.CHANGE

      const editableTemplate = this.$store.getters.getTemplateById(this.templateId);
      this.templateTitle = editableTemplate.Title
      this.queue = editableTemplate.Queue
    }
  }
}

const routes = [{
    path: "/",
    component: LoginPage
  },
  {
    path: "/main",
    component: MainPage
  },
  {
    name: "create_group",
    path: "/group/create",
    component: EditPage
  },
  {
    name: "create_template",
    path: "/group/:groupId/template/create",
    component: EditPage,
    props: true
  },
  {
    name: "edit_template",
    path: "/group/:groupId/template/:templateId/edit",
    component: EditPage,
    props: true
  }
]

const router = new VueRouter({
  routes
})

const app = new Vue({
  router,
  store
}).$mount("#app")
