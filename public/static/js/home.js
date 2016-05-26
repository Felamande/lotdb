window.onload = function() {

	Vue.use(VueResource)
	window.vm = new Vue({
		el: ".container",
		data: {
			errMsg: "",
			sum: 0,
			filters: [{}],
			maxFilterNum: 35,
			minFilterNum: 1,
			maxSum: 165,
			minSum: 15,
			results: [],
			deleted: [],
			firstRun: true,
			showDeleted: false,
			showFilters: false,
			debug: false,
			rePerPage: 25,
			curPageIdx: 0,
			selected: [],
			showTimeClose:true

		},
		components: {
			modal: VueStrap.modal,
			tabset: VueStrap.tabset,
			tab: VueStrap.tab,
			vselect: VueStrap.select,
			voption: VueStrap.option,
		},
		methods: {
			postFilter: function() {
				this.showFilters = false
				this.$http.post("/query", this.queryParams, { headers: { "Content-Type": "application/json" } }).then(
					function(response) {
						this.deleted = []
						data = response.data
						if (!data.success) {
							this.results = []
							this.errMsg = data.err
						} else if (data.re == null || data.re.length == 0) {
							this.results = []
							this.errMsg = "找不到满足条件的结果"
						} else {
							this.$set("results", data.re)
							this.errMsg = ""
						}
						this.firstRun = false

					}, function(error) { this.errMsg = error.data })

			},
			remResult: function(re) {
				this.results.$remove(re)
				// this.curPageSlice.$remove(re)
				this.deleted.push(re)
			},
			addDeletedResult: function(re) {
				this.deleted.$remove(re)
				this.results.push(re)
			},
			remFilter: function(f) {
				this.filters.$remove(f)
			},
			addFilter: function() {
				this.filters.push({})
			},
            delAllFilters:function(){
                this.filters = [{}]
            },
			postClickTooLate:function(){
				this.$http.post("/toolate", {"clicked":true}, { headers: { "Content-Type": "application/json" } }).then(
					function(response) {}, function(error) {})
				this.showTimeClose = false
			}
		},
		computed: {
			queryParams: function() {
				var data = {}
				data["sum"] = this.sum
				data["include"] = this.filterByType.include
                data["exclude"] = this.filterByType.exclude
				return data
			},
			maxPageIdx: function() {
				return parseInt(this.results.length / 25) + 1
			},
			filterByType: function() {
				var fbt = { exclude: [], include: [] }
                
				for (i in this.filters) {
                    filter = this.filters[i]
                    switch(filter.type){
                        case "include":
                            fbt.include.push(filter.value)
                        case "exclude":
                            for(i = filter.start;i<=filter.end;i++){
                                fbt.exclude.push(i)
                            }
                    }
				}
				return fbt
			},
			currPageSlice: function() {
				var next = this.curPageIdx + 1
				var start = this.curPageIdx * 25
				var end = next * 25
				return this.results.slice(start, end)
			},
			showTime:function(){
				var d = new Date()
				var h = d.getHours()
				var m = d.getMinutes()
				if (h<=4){
					if(h>=1&&m>=20){		
						return true
					}
				}
				return false
			},
			now:function(){
				var d = new Date()
				var h = d.getHours()
				var m = d.getMinutes()
				return h+":"+m
			}
			
		}
	})
}