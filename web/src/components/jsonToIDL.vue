<template>
  <div>
    <el-row :gutter="24">
      <el-form ref="elForm" :model="formData" :rules="rules" size="medium" label-width="100px"
               label-position="left">
        <el-col :span="24">
          <el-row>
          </el-row>
        </el-col>
        <el-col :span="12">
          <el-row>
            <el-col :span="24">
              <el-form-item label-width="0" prop="JSON">
                <el-input v-model="formData.JSON" type="textarea" placeholder="请输入待转化JSON"
                          :autosize="{minRows: 25, maxRows: 1000}" :style="{width: '100%'}"></el-input>
              </el-form-item>
            </el-col>
          </el-row>
        </el-col>
        <el-col :span="12">
          <el-row>
            <el-col :span="24">
              <el-form-item label-width="0" prop="IDL">
                <el-input v-model="formData.IDL" type="textarea" placeholder="输出IDL" readonly
                          :autosize="{minRows: 25, maxRows: 1000}" :style="{width: '100%'}"></el-input>
              </el-form-item>
            </el-col>
          </el-row>
        </el-col>
        <el-col :span="24">
          <el-row>
            <el-col :span="12">
              <el-row>
                <el-col :span="12">
                  <el-form-item label="示例" prop="example">
                    <el-select v-model="formData.example" placeholder="请选择示例" clearable
                               :style="{width: '100%'}" @change="selectChange">
                      <el-option v-for="(item, index) in exampleOptions" :key="index" :label="item.label"
                                 :value="item.value" :disabled="item.disabled"></el-option>
                    </el-select>
                  </el-form-item>
                </el-col>
              </el-row>
            </el-col>
            <!-- <el-col :span="4">
              <el-row>
              </el-row>
            </el-col> -->

            <el-col :span="12">
              <el-form-item size="large">
                <el-button @click="resetForm" type="info" icon="el-icon-milk-tea" size="medium">清空</el-button>
                <el-button type="primary" icon="el-icon-reading" size="medium" @click="submitForm">转化</el-button>
              </el-form-item>
            </el-col>

          </el-row>
        </el-col>
      </el-form>
    </el-row>
  </div>
</template>
<script>
export default {
  name: 'jsonToIDL',
  components: {},
  props: [],
  data() {
    return {
      formData: {
        JSON: undefined,
        IDL: undefined,
        example: "",
      },
      rules: {
        JSON: [{
          required: true,
          message: '请输入待转化JSON',
          trigger: 'blur'
        }],
        IDL: [{
          message: '输出IDL',
          trigger: 'blur'
        }],
        example: [{
          message: '请选择示例',
          trigger: 'change'
        }],
      },
      exampleOptions: [{
        "label": "例子1",
        "value": `{
	"name": "yangyang",
	"age": 25,
	"mail": false,
	"graduation": ["shaoyangxue", 1222]
}`
      }, {
        "label": "例子2",
        "value": `{
    "sites": [
    { "name":"github" , "url":"github.com/zy84338719" },
    { "name":"blog" , "url":"blog.murphyyi.com" }
    ]
}`
      }],
    }
  },
  computed: {},
  watch: {},
  created() {
  },
  mounted() {
  },
  methods: {
    submitForm() {
      this.$refs['elForm'].validate(valid => {
        if (!valid)
          return
        // TODO 提交表单
      })
      // eslint-disable-next-line no-console
      console.log(this.formData.JSON)
      fetch("api/json2IDL", {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({
          JSON: this.formData.JSON
        })
      })
      .then(
          function (data) {
            return data.json()
          }
      )
      .then(res => {
        // eslint-disable-next-line no-console
        console.log(res)
        if (res.code === 0) {
          this.formData.IDL = res.data.idl
        }else {
          this.formData.IDL = "出错了"
        }
      })
    },
    resetForm() {
      this.$refs['elForm'].resetFields()
    },
    selectChange(v) {
      this.formData.JSON = v
    }
  }
}

</script>
<style>
</style>
