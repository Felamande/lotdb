window.onload = function() {

	Vue.use(VueResource)

	Vue.transition('del', {})
	Vue.directive("highlight", function(value) {
		var filters = this.vm.$data.filters
		if (!filters) {
			return
		}
		for (i in filters) {
			if (value == filters[i].value) {
				Vue.util.addClass(this.el, "red")
			}
		}
	})
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
			selected: []

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
							this.setPageSlice(0)
							this.errMsg = ""
						}
						this.firstRun = false

					}, function(error) { this.errMsg = error.data })

			},
			remResult: function(re) {
				this.results.$remove(re)
				this.curPageSlice.$remove(re)
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

			highlight: function(n) {
				for (i in this.filters) {
					if (value == this.filters[i].value) {
						return "red"
					}
				}
			}

		},
		computed: {
			queryParams: function() {
				var data = {}
				data["sum"] = this.sum
				data["filters"] = this.filters
				return data
			},
			maxPageIdx: function() {
				return parseInt(this.results.length / 25) + 1
			},
			filterByType: function() {
				var fbt = { exclude: [], include: [] }
				for (i in this.filters) {


					if (this.filters[i].value == 0 || !this.filters[i].type) continue;
					fbt[this.filters[i].type].push(this.filters[i].value)
				}
				return fbt
			},
			currPageSlice: function() {
				var next = this.curPageIdx + 1
				var start = this.curPageIdx * 25
				var end = next * 25
				return this.results.slice(start, end)
			}
		}
	})
}