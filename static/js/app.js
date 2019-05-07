const api = {
  async connectToMq(url, user, pass) {
    const response = await axios.post('/api', {"jsonrpc": "2.0", "method": "MqService.Connect", "params": [{"URL": "tcp://localhost:61613"}], "id": 1});
    return response.data;
  },
  getAllGroups() {
    return new Promise((resolve, reject) => {
       setTimeout(() => {
          resolve([{
              title: 'group1',
              id: 1
            },
            {
              title: 'group2',
              id: 2
            },
            {
              title: 'group3',
              id: 3
            }
          ]);

    	}, 3000);
    })
  } 
}

const store = new Vuex.Store({
  state: {
    groups: [],
    templates: [{
        title: 'group1_template1',
        queue: 'group1_template1_queue',
        id: 1,
        groupId: 1
      },
      {
        title: 'group1_template2',
        queue: 'group1_template2_queue',
        id: 2,
        groupId: 1
      },
      {
        title: 'group2_template1',
        queue: 'group2_template1_queue',
        id: 3,
        groupId: 2
      },
      {
        title: 'group3_template1',
        queue: 'group3_template1_queue',
        id: 4,
        groupId: 3
      }
    ]
  },
  mutations: {
    addGroup (state, group) {
      state.groups.push(group);
    }
  },
  actions: {
    async signIn (context) {
      const res = await api.connectToMq();
      console.log(res)
      //TODO: check res.error
      const groups = await api.getAllGroups();
      groups.forEach((group) => {
        context.commit('addGroup', group)
      });
    }
  },
  getters: {
    getTemplatesByGroupId: (state) => (id) => {
      return state.templates.filter(template => template.groupId === id);
    },
    getTemplateById: (state) => (id) => {
      return state.templates.find(template => template.groupId === id);
    },
    getGroupById: (state) => (id) => {
      return state.groups.find(group => group.id === id);
    }
  }
})

const Login = {
  template: `
  	<form @submit.prevent="onSubmit" class="pure-form pure-form-aligned">
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
          <input id="pass" v-model="pass">
        </div>
        <div class="pure-controls">
          <button type="submit" :disabled="!isSubmitButtonEnabled" class="pure-button pure-button-primary">Sign in</button>
        </div>
      </fieldset>
		</form>
	`,
  methods: {
    onSubmit: function() {
      //console.log({ name: this.name, email: this.email });
      //alert('Connecting. URL: ' + this.url + ', user: ' + this.user + ', pass: ' + this.pass);
      this.isSubmitButtonEnabled = false;

      this.$store.dispatch('signIn').then(() => {
        this.$router.push("/main");
      })
      
    }
  },
  data: function() {
    return {
      url: '',
      user: '',
      pass: '',
      isSubmitButtonEnabled: true
    }
  }
}

const Main = {
  template: `
  	<form @submit.prevent="onSubmit" class="pure-form pure-form-aligned">
    	<fieldset>
        <div class="pure-control-group">
          <label for="group">Group:</label>
          <select id="group" v-model="selectedGroupId">
            <option v-for="group in groupList" v-bind:value="group.id">
              {{ group.title }}
            </option>
          </select>
        </div>
        <div class="pure-control-group">
          <label for="template">Template:</label>
          <select id="template" v-model="selectedTemplateId">
            <option v-for="template in templateList" v-bind:value="template.id">
              {{ template.title }}
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
          <textarea id="message" v-model="message" style="vertical-align: top;" />
        </div>
        <div class="pure-controls">
          <button type="submit" class="pure-button pure-button-primary">Send</button>
        </div>
      </fieldset>
    </form>
	`,
  methods: {
    onSubmit: function() {
      alert('Not implemented yet!');
    },
    onEdit: function() {
      this.$router.push({
        name: 'edit_template',
        params: {
          groupId: this.selectedGroupId,
          templateId: this.selectedTemplateId
        }
      });
    }
  },
  data: function() {
    return {
      selectedGroupId: 0,
      selectedTemplateId: 0,
      message: ''
    }
  },
  computed: {
    selectedTemplateQueue() {
    	const selectedTemplate = this.$store.getters.getTemplateById(this.selectedTemplateId);
      return selectedTemplate ? selectedTemplate.queue : '';
    },
    groupList() {
      return this.$store.state.groups
    },
    templateList: function() {
      this.selectedTemplateId = 0;

      return this.$store.getters.getTemplatesByGroupId(this.selectedGroupId);
    }
  }
}

const EditTemplate = {
  props: ['groupId', 'templateId'],
  template: `
  	<form class="pure-form pure-form-aligned">
    	<fieldset>
        <div class="pure-control-group">
          <label for="groupTitle">Group title:</label>
          <input id="groupTitle" v-model="groupTitle">
				</div>
        <div class="pure-control-group">
          <label for="templateTitle">Template title:</label>
          <input id="templateTitle" v-model="templateTitle">
        </div>
        <div class="pure-control-group">
          <label for="queue">Queue:</label>
          <input id="queue" v-model="queue">
        </div>
        <div class="pure-controls">
        	<button @click.prevent="onSave" class="pure-button pure-button-primary">Save</button>
          <button @click.prevent="onCancel" class="pure-button">Cancel</button>
        </div>
      </fieldset>
    </form>
	`,
  methods: {
    onSave: function() {
      alert('Not implemented yet!');
    },
    onCancel: function() {
      this.$router.push("/main");
    }
  },
  data: function() {
    return {
      groupTitle: '',
      templateTitle: '',
      queue: ''
    }
  },
  created: function() {
  	const editableTemplate = this.$store.getters.getTemplateById(this.templateId);
    const editableGroup = this.$store.getters.getGroupById(this.groupId);
    this.groupTitle = editableGroup.title
    this.templateTitle = editableTemplate.title
    this.queue = editableTemplate.queue
  }
}

const routes = [{
    path: '/',
    component: Login
  },
  {
    path: '/main',
    component: Main
  },
  {
    name: 'edit_template',
    path: '/group/:groupId/template/:templateId/edit',
    component: EditTemplate,
    props: true
  }
]

const router = new VueRouter({
  routes
})

const app = new Vue({
  router,
  store
}).$mount('#app')
