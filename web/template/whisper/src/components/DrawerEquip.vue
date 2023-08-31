<script setup>

</script>
<script>
import axios from 'axios';

export default {
  data() {
    return {
      sideDrawer: false,
      sideDrawerData: {
        "current": {},
        "from": [],
        "into": [],
        "gapPriceFrom": 0,
      },
    }
  },
  watch: {
    'formData.keywords'() {
      // 使用 Axios 发起请求获取服务器数据
      if (this.formData.keywords.length === 0) {
        return
      }
      axios.post('/equip/filter', this.formData)
          .then(response => {
            this.result = response.data.data
          })
          .catch(error => {
            console.error('Error fetching server data:', error);
          });
    },
    'formData.platform'(n) {

      this.platformCN = n === '0' ? '端游' : '手游'
      // 使用 Axios 发起请求获取服务器数据
      if (this.formData.keywords.length === 0) {
        return
      }
      axios.post('/equip/filter', this.formData)
          .then(response => {
            this.result = response.data.data
          })
          .catch(error => {
            console.error('Error fetching server data:', error);
          });
    },
  },
  methods: {
    showDrawer(platform, version, id) {
      axios.post('/equip/roadmap', {
        'platform': platform,
        'version': version,
        'id': id,
      }).then(response => {
        this.sideDrawer = true
        this.sideDrawerData = response.data.data
      }).catch(error => {
        console.error('Error fetching server data:', error);
      });
    }
  },
  created() {
  },
  computed: {},
  mounted() {}
}
</script>
<a-drawer
    v-model:open="sideDrawer"
    class="custom-class"
    root-class-name="root-class-name"
    title="合成路线"
    placement="right"
>
<div class="equip-roadmap equip-into">
  <template v-for="(equip ,index) in sideDrawerData['into']" :key="index">
    <a-popover placement="bottom" arrow-point-at-center>
      <template #content>
        <div class="roadmap-item">
                    <span class="roadmap-item-title">
                      {{ equip.name }}
                    </span>
          <a-tag>
                    <span class="roadmap-item-price">
                      价格:{{ equip.price }}
                    </span>
          </a-tag>
        </div>
        <span v-html="equip.desc"></span>
      </template>
      <img :src="equip.icon">
    </a-popover>
  </template>
</div>
<a-divider>
  <a-popover placement="bottom" arrow-point-at-center>
    <template #content>
      <div class="roadmap-item">
                    <span class="roadmap-item-title">
                      {{ sideDrawerData['current'].name }}
                    </span>
        <a-tag>
                    <span class="roadmap-item-price">
                      价格:{{ sideDrawerData['current'].price }}
                    </span>
        </a-tag>
      </div>
      <span v-html="sideDrawerData['current'].desc"></span>
    </template>
    <img class="equip-roadmap equip-current" :src="sideDrawerData['current'].icon" alt="">
  </a-popover>

</a-divider>
<div class="equip-roadmap equip-from">
  <template v-for="(equip ,index) in sideDrawerData['from']" :key="index">
    <a-popover placement="bottom" arrow-point-at-center>
      <template #content>
        <div class="roadmap-item">
                    <span class="roadmap-item-title">
                      {{ equip.name }}
                    </span>
          <a-tag>
                    <span class="roadmap-item-price">
                      价格:{{ equip.price }}
                    </span>
          </a-tag>
        </div>
        <span v-html="equip.desc"></span>
      </template>
      <img :src="equip.icon">
    </a-popover>
  </template>
</div>

<table class="roadmap-detail">
  <tr v-for="(equip ,index) in sideDrawerData['from']" :key="index">
    <td>
      <img :src="equip.icon" alt="">
      <span>{{ equip.name }}</span>
      <span>, 价格: <em>{{ equip.price }}</em> </span>
    </td>
  </tr>
  <tr>
    <td>
      总价: <em>{{ sideDrawerData['current'].price }}</em> , 补差价: <em>{{ sideDrawerData.gapPriceFrom }}</em>
    </td>
  </tr>
</table>
</a-drawer>