<script setup>

</script>
<script>
import axios from 'axios';
import LoadingForm from "@/components/LoadingForm.vue";
import LoadingList from "@/components/LoadingList.vue";
import ListEquip from "@/components/ListEquip.vue";

export default {
  emits: ['loadingEvent'],
  components: {
    LoadingForm, LoadingList, ListEquip
  },
  data() {
    return {
      showCheckbox: true,
      loadingState: {
        condCheckbox: true,
      },
      formData: {
        platform: '0',
        keywords: [],
        equipTypes: [
          {
            cate: '',
            sub_cate: [{
              name: '',
              keywordsStr: '',
              keywordsSlice: [],
            }],
          }
        ],
      },
      SkeletonState: {
        show: false,
        isLoading: false,
      },
      query: {
        tips: '',
        referer: '',
        data: {},
      }
    }
  },
  watch: {
    formData: {
      handler(newValue) {
        this.$emit('loadingEvent', 0)
        this.platformCN = newValue === '0' ? '端游' : '手游'
        // 使用 Axios 发起请求获取服务器数据
        if (this.formData.keywords.length === 0) {
          return
        }

        this.SkeletonState.show = true
        this.$emit('loadingEvent', 30)
        axios.post('/equip/filter', this.formData)
            .then(response => {
              this.query.data = response.data.data.list
              this.query.tips = response.data.data.tips
              this.query.referer = "equip-box"
            })
            .catch(error => {
              console.error('Error fetching server data:', error);
            }).finally(() => {
              this.SkeletonState.show = false
              this.$emit('loadingEvent', 100)
            }
        );
      },
      deep: true,
    },
  },
  methods: {},
  created() {
  },
  computed: {},
  mounted() {
    // 使用 Axios 发起请求获取服务器数据
    axios.get('/equip/types')
        .then(response => {
              // 将服务器返回的数据更新到组件的 serverData 字段
              this.formData.equipTypes = response.data.data.types;
            }
        ).catch(error => {
          console.error('Error fetching server data:', error);
        }
    ).finally(() => {
          this.loadingState.condCheckbox = false;
        }
    );
  }
}
</script>

<template>
  <a-space direction="vertical" :style="{ width: '100%' }" class="wrap">
    <a-layout>
      <a-layout-content>
        <a-form :model="formData" name="equip-box">
          <div class="platform-switch">
            <a-radio-group v-model:value="formData.platform" name="radioGroup">
              <a-radio value=0>端游</a-radio>
              <a-radio value=1>手游</a-radio>
            </a-radio-group>
            <a-button type="link" class="cond-drawer" size="small" @click="showCheckbox = !showCheckbox">
              {{ showCheckbox ? '收起条件' : '展开条件' }}
            </a-button>
          </div>

          <LoadingForm v-if="loadingState.condCheckbox"/>
          <Transition>
            <template v-if="showCheckbox && !loadingState.condCheckbox">
              <a-checkbox-group v-model:value="formData.keywords">
                <a-descriptions v-for="(cate, index) in formData.equipTypes" :title="cate.cate" :key="index">
                  <a-descriptions-item>
                    <a-checkbox v-for="(sub_cate, ii) in cate.sub_cate" :key="ii" :id="index+'-'+ii+'-equip-checkbox'"
                                :value="sub_cate.keywordsStr">
                      <a-tooltip placement="top" class="equip-box-tooltip">
                        <template #title>
                          <span>{{ sub_cate.keywordsStr }}</span>
                        </template>
                        {{ sub_cate.name }}
                      </a-tooltip>
                    </a-checkbox>
                  </a-descriptions-item>
                </a-descriptions>
              </a-checkbox-group>
            </template>
          </Transition>
        </a-form>

        <LoadingList
            v-if="SkeletonState.show"
            :skeleton-state="SkeletonState"/>

        <ListEquip
            v-if="!SkeletonState.show"
            :query-result="query"
        />

        <a-back-top/>
      </a-layout-content>
    </a-layout>
  </a-space>
</template>

<style scoped>
.v-enter-active,
.v-leave-active {
  transition: opacity 0.5s ease;
}

.v-enter-from,
.v-leave-to {
  opacity: 0;
}
</style>
