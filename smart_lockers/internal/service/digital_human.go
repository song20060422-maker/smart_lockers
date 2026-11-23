// @ticket 07-数字人Coze：本地演示问答引擎——读取项目真实 faq.json / brand.md，内置关键词匹配，不接 Coze、无需任何凭证
// 定义包 service：为 HTTP / WebSocket handler 提供可复用的业务能力

package service

import (
	// 引入 encoding/json：解析真实 faq.json 并序列化回复对象

	"encoding/json"
	// 引入 os：读取知识文件与读取环境变量覆盖路径

	"os"
	// 引入 path/filepath：拼接跨平台的知识文件路径

	"path/filepath"
	// 引入 strings：做文本归一化、关键词命中与品牌文档解析

	"strings"
	// 引入 sync：保证知识库单例只加载一次

	"sync"
	// 引入 unicode：过滤标点与空白并统一大小写

	"unicode"
)

// FAQItem 对应 knowledge/faq.json 中 items 的单个条目，字段与真实 JSON 完全一致
// 定义结构体 FAQItem：用于承载一条真实 FAQ 知识

type FAQItem struct {
	// 结构体字段 ID，类型为 string，附带标签 `json:"id"`

	ID string `json:"id"`
	// 结构体字段 Tags，类型为 []string，附带标签 `json:"tags"`

	Tags []string `json:"tags"`
	// 结构体字段 Weight，类型为 float64，附带标签 `json:"weight"`

	Weight float64 `json:"weight"`
	// 结构体字段 Q，类型为 string，附带标签 `json:"q"`

	Q string `json:"q"`
	// 结构体字段 A，类型为 string，附带标签 `json:"a"`

	A string `json:"a"`
	// 当前代码块结束

}

// faqDocument 对应 knowledge/faq.json 根结构
// 定义结构体 faqDocument：用于反序列化真实 FAQ 文件

type faqDocument struct {
	// 结构体字段 Ticket，类型为 string，附带标签 `json:"_ticket"`

	Ticket string `json:"_ticket"`
	// 结构体字段 Items，类型为 []FAQItem，附带标签 `json:"items"`

	Items []FAQItem `json:"items"`
	// 当前代码块结束

}

// Live2DAction 严格对齐数字人前端 src/types/index.ts 中 Live2DAction 的真实取值
// 定义结构体 Live2DAction：用于描述数字人表情/动作联动

type Live2DAction struct {
	// 结构体字段 Action，取值为 nod/smile/shake/blink/transform/idle，附带标签 `json:"action"`

	Action string `json:"action"`
	// 结构体字段 Emotion，类型为 string，附带标签 `json:"emotion,omitempty"`

	Emotion string `json:"emotion,omitempty"`
	// 结构体字段 DurationMS，前端 emotionAction 默认 1200，附带标签 `json:"duration_ms,omitempty"`

	DurationMS int `json:"duration_ms,omitempty"`
	// 当前代码块结束

}

// AgentReplyPayload 严格对齐数字人前端 src/types/index.ts 中 AgentReplyPayload，不增加臆造字段
// 定义结构体 AgentReplyPayload：用于承载智能体回复文本与情绪/场景/动作

type AgentReplyPayload struct {
	// 结构体字段 Text，类型为 string，附带标签 `json:"text"`

	Text string `json:"text"`
	// 结构体字段 Emotion，前端 emotionAction 支持 happy/positive/sad/neutral，附带标签 `json:"emotion"`

	Emotion string `json:"emotion"`
	// 结构体字段 Scene，前端默认 chat，附带标签 `json:"scene"`

	Scene string `json:"scene"`
	// 结构体字段 Action，类型为 *Live2DAction，附带标签 `json:"action,omitempty"`

	Action *Live2DAction `json:"action,omitempty"`
	// 当前代码块结束

}

// AgentReply 严格对齐数字人前端 WebSocketClient.ts 仅识别的 assistant 消息结构
// 定义结构体 AgentReply：用于序列化成单条 JSON 文本帧返回给前端

type AgentReply struct {
	// 结构体字段 Role，固定为 assistant，附带标签 `json:"role"`

	Role string `json:"role"`
	// 结构体字段 Payload，类型为 AgentReplyPayload，附带标签 `json:"payload"`

	Payload AgentReplyPayload `json:"payload"`
	// 当前代码块结束

}

// localKnowledge 保存从项目真实文件加载的知识内容
// 定义结构体 localKnowledge：用于缓存 FAQ 条目与品牌介绍

type localKnowledge struct {
	// 结构体字段 FAQs，类型为 []FAQItem，来自 faq.json 的 items

	FAQs []FAQItem
	// 结构体字段 BrandIntro，来自 brand.md 正文首段

	BrandIntro string
	// 结构体字段 BrandLines，来自 brand.md 的列表条目

	BrandLines []string
	// 当前代码块结束

}

// RoleProfile 角色人设档案：来自 knowledge/roles.json，实现三角色内容互不相同 @ticket 07
// 定义结构体 RoleProfile：承载单角色的人设、欢迎语与专属 FAQ 文件
type RoleProfile struct {
	// 结构体字段 ID，角色标识 sender/receiver/service，附带标签 `json:"id"`
	ID string `json:"id"`
	// 结构体字段 Name，角色中文名，附带标签 `json:"name"`
	Name string `json:"name"`
	// 结构体字段 Persona，角色人设话术，附带标签 `json:"persona"`
	Persona string `json:"persona"`
	// 结构体字段 Welcome，角色欢迎语，附带标签 `json:"welcome"`
	Welcome string `json:"welcome"`
	// 结构体字段 FAQFile，角色专属 FAQ 文件名，附带标签 `json:"faqFile"`
	FAQFile string `json:"faqFile"`
	// 当前代码块结束

}

// rolesDocument 对应 knowledge/roles.json 根结构
// 定义结构体 rolesDocument：用于反序列化角色档案配置
type rolesDocument struct {
	// 结构体字段 Roles，类型为 []RoleProfile，附带标签 `json:"roles"`
	Roles []RoleProfile `json:"roles"`
	// 当前代码块结束

}

// roleKnowledge 单个角色的完整运行态知识：人设 + 欢迎语 + 专属 FAQ
// 定义结构体 roleKnowledge：合并角色档案与专属 FAQ，供 BuildReply 按角色取用
type roleKnowledge struct {
	// 结构体字段 Profile，角色人设档案
	Profile RoleProfile
	// 结构体字段 FAQs，角色专属 FAQ 条目
	FAQs []FAQItem
	// 当前代码块结束

}

// 开始定义一组包级/局部变量的 var 声明块

var (
	// 执行 loadOnce sync.Once，保证知识库只初始化一次

	loadOnce sync.Once
	// 执行 knowledge *localKnowledge，保存单例知识库

	knowledge *localKnowledge
	// 执行 roleLoadOnce sync.Once，保证角色知识只初始化一次 @ticket 07

	roleLoadOnce sync.Once
	// 执行 roleKnowledgeMap map[string]*roleKnowledge，按角色 id 缓存专属知识 @ticket 07

	roleKnowledgeMap map[string]*roleKnowledge
)

// embeddedFAQs 是 faq.json 中 q01-q05 的逐字副本，仅在文件读取失败时兜底，内容非臆造
// 定义函数 embeddedFAQs：返回内置的真实 FAQ 条目

func embeddedFAQs() []FAQItem {
	// 返回包含五条真实 FAQ 的切片，逐字来自 knowledge/faq.json

	return []FAQItem{
		// 执行 q01 寄存条目

		{ID: "q01", Tags: []string{"寄存"}, Weight: 1.0, Q: "如何寄存物品？", A: "在首页选择附近寄存点，选择空闲柜机后扫码支付即可寄存，寄存成功后会收到开柜凭证。"},
		// 执行 q02 收费条目

		{ID: "q02", Tags: []string{"收费"}, Weight: 1.0, Q: "寄存怎么收费？", A: "按柜型与寄存时长计费，小柜起步价更低，具体费用以柜机屏幕与订单页展示为准。"},
		// 执行 q03 取件条目

		{ID: "q03", Tags: []string{"取件"}, Weight: 1.0, Q: "怎么取出物品？", A: "在订单详情点击临时开柜，或用凭证码在柜机屏幕取件；超时取件需先补交费用。"},
		// 执行 q04 故障条目

		{ID: "q04", Tags: []string{"故障"}, Weight: 0.9, Q: "柜门打不开怎么办？", A: "请先确认订单状态与凭证码正确；若仍无法开启，可联系在线客服或拨打柜机上的服务电话。"},
		// 执行 q05 场景条目

		{ID: "q05", Tags: []string{"品牌", "景区"}, Weight: 0.8, Q: "你们支持哪些场景？", A: "智能快寄柜覆盖景区、商圈、校园、交通枢纽等场景，提供自助寄存与智能客服服务。"},
		// 当前代码块结束

	}
	// 当前代码块结束

}

// embeddedBrandIntro 是 brand.md 正文首段的逐字副本，仅在文件读取失败时兜底
// 定义常量 embeddedBrandIntro：保存真实品牌介绍句

const embeddedBrandIntro = "智能快寄柜是智能自助寄存系统，提供就近寄存、扫码开柜、订单管理与智能客服一体化服务。"

// findKnowledgeDir 按优先级探测真实知识目录，支持环境变量覆盖
// 定义函数 findKnowledgeDir：返回包含 faq.json 与 brand.md 的目录

func findKnowledgeDir() string {
	// 读取环境变量 DIGITAL_HUMAN_KB_DIR 作为最高优先级目录

	dir := os.Getenv("DIGITAL_HUMAN_KB_DIR")
	// 条件判断：dir 非空且存在 faq.json 时执行下面逻辑

	if dir != "" && fileExists(filepath.Join(dir, "faq.json")) {
		// 返回 dir 给调用方，结束当前函数

		return dir
		// 当前代码块结束

	}
	// 定义候选目录列表：覆盖从项目根启动与从后端子目录启动两种工作目录

	candidates := []string{
		// 从 smart_lockers/smart_lockers 启动时的相对路径

		filepath.Join("..", "digital_human", "coze", "knowledge"),
		// 从 smart_lockers 项目根启动时的相对路径

		filepath.Join("digital_human", "coze", "knowledge"),
		// 从更深层级启动时的相对路径

		filepath.Join("..", "..", "digital_human", "coze", "knowledge"),
	}
	// 遍历候选目录并校验 faq.json 是否存在

	for _, candidate := range candidates {
		// 条件判断：fileExists(filepath.Join(candidate, "faq.json")) 时执行下面逻辑

		if fileExists(filepath.Join(candidate, "faq.json")) {
			// 返回 candidate 给调用方，结束当前函数

			return candidate
			// 当前代码块结束

		}
		// 当前代码块结束

	}
	// 返回默认第一个候选目录，由后续读取失败走内嵌兜底

	return candidates[0]
	// 当前代码块结束

}

// fileExists 判断给定路径是否为可读文件
// 定义函数 fileExists：封装文件存在性检查

func fileExists(path string) bool {
	// 执行 info, err := os.Stat(path)，获取文件信息

	info, err := os.Stat(path)
	// 条件判断：err != nil || info.IsDir() 时执行下面逻辑

	if err != nil || info.IsDir() {
		// 返回 false 给调用方，结束当前函数

		return false
		// 当前代码块结束

	}
	// 返回 true 给调用方，结束当前函数

	return true
	// 当前代码块结束

}

// parseBrand 从 brand.md Markdown 中提取正文首段与列表条目，去掉注释/标题/列表符号
// 定义函数 parseBrand：返回品牌介绍句与要点行

func parseBrand(md string) (string, []string) {
	// 定义并初始化局部变量 intro，保存正文首段

	intro := ""
	// 定义并初始化局部变量 lines，保存要点行

	var lines []string
	// 按换行拆分 Markdown 文档

	for _, raw := range strings.Split(md, "\n") {
		// 去除首尾空白得到 line

		line := strings.TrimSpace(raw)
		// 条件判断：跳过空行、标题行、HTML 注释行

		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, "<!--") {
			// 进入下一轮循环

			continue
			// 当前代码块结束

		}
		// 条件判断：处理列表条目行

		if strings.HasPrefix(line, "- ") {
			// 去掉列表符号后追加到 lines

			lines = append(lines, strings.TrimSpace(strings.TrimPrefix(line, "- ")))
			// 进入下一轮循环

			continue
			// 当前代码块结束

		}
		// 条件判断：尚未设置正文首段时执行下面逻辑

		if intro == "" {
			// 将当前行作为正文首段

			intro = line
			// 当前代码块结束

		}
		// 当前代码块结束

	}
	// 返回 intro 与 lines 给调用方，结束当前函数

	return intro, lines
	// 当前代码块结束

}

// getKnowledge 单例加载真实 faq.json 与 brand.md，读取失败时用逐字内嵌副本兜底
// 定义函数 getKnowledge：返回只读知识库实例

func getKnowledge() *localKnowledge {
	// 执行 loadOnce.Do 保证并发安全且只加载一次

	loadOnce.Do(func() {
		// 定义并初始化局部变量 k

		k := &localKnowledge{}
		// 探测知识目录

		dir := findKnowledgeDir()
		// 读取真实 faq.json 文件

		if b, err := os.ReadFile(filepath.Join(dir, "faq.json")); err == nil {
			// 定义并初始化局部变量 doc

			var doc faqDocument
			// 条件判断：解析 JSON 成功后执行下面逻辑

			if json.Unmarshal(b, &doc) == nil {
				// 为 k.FAQs 赋值 doc.Items

				k.FAQs = doc.Items
				// 当前代码块结束

			}
			// 当前代码块结束

		}
		// 条件判断：FAQ 为空时使用内嵌真实条目兜底

		if len(k.FAQs) == 0 {
			// 为 k.FAQs 赋值 embeddedFAQs()

			k.FAQs = embeddedFAQs()
			// 当前代码块结束

		}
		// 读取真实 brand.md 文件

		if b, err := os.ReadFile(filepath.Join(dir, "brand.md")); err == nil {
			// 解析品牌文档得到正文首段与要点行

			k.BrandIntro, k.BrandLines = parseBrand(string(b))
			// 当前代码块结束

		}
		// 条件判断：品牌介绍为空时使用内嵌真实文案兜底

		if k.BrandIntro == "" {
			// 为 k.BrandIntro 赋值 embeddedBrandIntro

			k.BrandIntro = embeddedBrandIntro
			// 当前代码块结束

		}
		// 为全局单例 knowledge 赋值

		knowledge = k
		// 执行 })，完成一次性初始化

	})
	// 返回 knowledge 给调用方，结束当前函数

	return knowledge
	// 当前代码块结束

}

// embeddedRoleProfiles 是 roles.json 三角色人设的逐字内嵌副本，仅在文件读取失败时兜底，内容来自真实配置 @ticket 07
// 定义函数 embeddedRoleProfiles：返回内置的三角色档案与专属 FAQ 文件名
func embeddedRoleProfiles() []RoleProfile {
	// 返回三角色档案切片，逐字来自 knowledge/roles.json
	return []RoleProfile{
		{ID: "sender", Name: "寄件引导员", Persona: "您好，我是寄件引导员小寄，专注帮您完成寄存下单、包装建议与开柜指引。", Welcome: "您好！我是寄件引导员小寄～从选柜、包装到支付开柜，我全程陪您完成寄件，请问您想寄存什么物品呢？", FAQFile: "faq.sender.json"},
		{ID: "receiver", Name: "取件指引员", Persona: "您好，我是取件指引员小取，专注帮您凭码取件、处理超时补费与代取咨询。", Welcome: "您好！我是取件指引员小取～凭取件码或订单二维码就能开柜，超时补费、代取问题也都能找我，请问您的取件码是多少？", FAQFile: "faq.receiver.json"},
		{ID: "service", Name: "智能客服专员", Persona: "您好，我是智能客服专员小服，负责品牌咨询、故障报修与人工客服转接。", Welcome: "您好！我是智能客服专员小服～品牌介绍、故障报修、投诉建议和人工客服转接都由我负责，请问需要哪方面帮助？", FAQFile: "faq.service.json"},
	}
}

// embeddedRoleFAQs 是三角色专属 FAQ 的内嵌逐字副本，仅在文件读取失败时兜底 @ticket 07
// 定义函数 embeddedRoleFAQs：按角色 id 返回内嵌 FAQ 条目
func embeddedRoleFAQs(id string) []FAQItem {
	// 按角色 id 分支返回对应内嵌 FAQ
	switch id {
	case "sender":
		return []FAQItem{
			{ID: "s01", Tags: []string{"寄件", "下单", "寄存"}, Weight: 1.0, Q: "如何寄件下单？", A: "寄件引导流程：在首页选择附近柜机，点击寄件，选择柜型后扫码支付，开门放入包裹并关好柜门即可，系统会自动推送开柜凭证给您。"},
			{ID: "s02", Tags: []string{"柜型", "尺寸", "大件"}, Weight: 1.0, Q: "寄件怎么选柜型？", A: "请根据物品体积选择：小柜适合文件、钥匙等小件，中柜适合衣物与日常包裹，大柜适合行李箱等大件；以柜机屏幕实时显示的空闲柜型为准。"},
			{ID: "s03", Tags: []string{"包装", "易碎", "防摔"}, Weight: 1.0, Q: "寄件需要自己包装吗？", A: "建议自行妥善包装，易碎品请加气泡膜或填充物固定，液体请密封防漏；包装完好能避免运输途中损坏，柜机内也提供应急打包提示。"},
			{ID: "s04", Tags: []string{"收费", "费用", "支付"}, Weight: 1.0, Q: "寄件怎么收费？", A: "寄件费用按柜型与寄存时长计算，小柜起步价更低，具体金额在下单页与柜机屏幕实时展示，支付成功后才会开柜。"},
			{ID: "s05", Tags: []string{"凭证", "开柜", "订单"}, Weight: 0.9, Q: "寄件后怎么拿到凭证？", A: "完成寄存后，凭证会通过订单页与短信推送给您，含取件码与开柜二维码，收件方凭此即可到柜机取件。"},
		}
	case "receiver":
		return []FAQItem{
			{ID: "r01", Tags: []string{"取件", "开柜", "凭证"}, Weight: 1.0, Q: "怎么取出物品？", A: "取件指引：在订单详情点击临时开柜，或在柜机屏幕选择取件并输入取件码/扫描二维码，柜门弹开后请及时取出并关好柜门。"},
			{ID: "r02", Tags: []string{"取件码", "凭证码", "在哪"}, Weight: 1.0, Q: "取件码在哪里看？", A: "取件码会在寄件成功的短信和订单详情页展示；若找不到，可让寄件方在订单详情查看并转发给您。"},
			{ID: "r03", Tags: []string{"超时", "补费", "过期"}, Weight: 1.0, Q: "超时取件要补费吗？", A: "免费寄存时长以订单约定为准，超出后需按柜型补交超时费用，补费完成后柜门才会开启，标准以柜机屏幕提示为准。"},
			{ID: "r04", Tags: []string{"代取", "帮别人", "转赠"}, Weight: 0.9, Q: "可以让别人代取吗？", A: "可以，把取件码或开柜二维码转发给对方即可代取；如支持转赠，也可在订单详情操作取件人变更。"},
			{ID: "r05", Tags: []string{"故障", "打不开", "柜门"}, Weight: 0.9, Q: "柜门打不开怎么办？", A: "请先确认取件码正确且订单状态正常；若仍无法开启，可联系在线客服或拨打柜机上的服务电话，我们会远程协助开柜。"},
		}
	default:
		return []FAQItem{
			{ID: "v01", Tags: []string{"品牌", "介绍", "你们是谁", "八维", "简介"}, Weight: 1.0, Q: "你们是什么品牌？", A: "我们是智能快寄柜，提供自助寄存与智能客服服务，覆盖景区、商圈、校园、交通枢纽等场景，随时为您解答寄存相关问题。"},
			{ID: "v02", Tags: []string{"场景", "覆盖", "景区", "校园"}, Weight: 0.9, Q: "支持哪些场景？", A: "智能快寄柜覆盖景区、商圈、校园、交通枢纽等多种场景，提供自助寄存、取件引导与智能客服服务。"},
			{ID: "v03", Tags: []string{"故障", "报修", "打不开", "异常"}, Weight: 1.0, Q: "设备故障怎么办？", A: "遇到柜门无法开启、扣费异常等问题，请先核对订单状态；若仍未解决，可联系在线客服或拨打柜机上的服务电话报修。"},
			{ID: "v04", Tags: []string{"人工", "客服", "转接", "电话"}, Weight: 1.0, Q: "怎么联系人工客服？", A: "您可以点击页面上的在线客服入口转接人工，或拨打柜机屏幕与凭证短信上的服务电话，工作时间内会有专员为您处理。"},
			{ID: "v05", Tags: []string{"投诉", "建议", "反馈", "评价"}, Weight: 0.9, Q: "如何投诉或提建议？", A: "欢迎通过在线客服提交投诉或建议，请描述问题与时间地点，我们会尽快跟进并反馈处理结果。"},
		}
	}
}

// getRoleKnowledge 单例加载 roles.json 与各角色专属 FAQ，按角色 id 提供差异化知识 @ticket 07
// 定义函数 getRoleKnowledge：返回按角色 id 索引的运行态知识 map
func getRoleKnowledge() map[string]*roleKnowledge {
	// 执行 roleLoadOnce.Do 保证并发安全且只加载一次
	roleLoadOnce.Do(func() {
		// 初始化 map
		roleKnowledgeMap = make(map[string]*roleKnowledge)
		// 探测知识目录
		dir := findKnowledgeDir()
		// 定义角色档案列表，默认使用内嵌副本
		profiles := embeddedRoleProfiles()
		// 读取真实 roles.json 文件
		if b, err := os.ReadFile(filepath.Join(dir, "roles.json")); err == nil {
			var doc rolesDocument
			if json.Unmarshal(b, &doc) == nil && len(doc.Roles) > 0 {
				profiles = doc.Roles
			}
		}
		// 遍历角色档案，逐个加载专属 FAQ
		for _, p := range profiles {
			rk := &roleKnowledge{Profile: p}
			// 读取角色专属 FAQ 文件
			if p.FAQFile != "" {
				if b, err := os.ReadFile(filepath.Join(dir, p.FAQFile)); err == nil {
					var doc faqDocument
					if json.Unmarshal(b, &doc) == nil {
						rk.FAQs = doc.Items
					}
				}
			}
			// FAQ 为空时用内嵌副本兜底
			if len(rk.FAQs) == 0 {
				rk.FAQs = embeddedRoleFAQs(p.ID)
			}
			roleKnowledgeMap[p.ID] = rk
		}
	})
	return roleKnowledgeMap
}

// welcomeForRole 返回指定角色的欢迎语，未识别角色回退到智能客服专员 @ticket 07
// 定义函数 welcomeForRole：供连接建立时主动推送角色欢迎语
func welcomeForRole(role string) string {
	// 取角色知识
	rk := getRoleKnowledge()[role]
	// 条件判断：命中角色时返回其欢迎语
	if rk != nil {
		return rk.Profile.Welcome
	}
	// 未命中时回退 service 角色欢迎语
	if fallback := getRoleKnowledge()["service"]; fallback != nil {
		return fallback.Profile.Welcome
	}
	return ""
}

// BuildWelcome 构造指定角色连接建立时的欢迎语回复，协议与 BuildReply 一致 @ticket 07
// 定义函数 BuildWelcome：返回可直接序列化推送的 assistant 欢迎语，走同一口型/动画链路
func BuildWelcome(role string) AgentReply {
	return AgentReply{
		Role: "assistant",
		Payload: AgentReplyPayload{
			Text:    welcomeForRole(role),
			Emotion: "happy",
			Scene:   "chat",
			Action: &Live2DAction{
				Action:     "smile",
				Emotion:    "happy",
				DurationMS: 1200,
			},
		},
	}
}

// normalize 归一化文本：过滤标点与空白并转小写，便于中文问答匹配
// 定义函数 normalize：返回可用于匹配的干净文本

func normalize(text string) string {
	// 定义并初始化字符串构造器 b

	var b strings.Builder
	// 遍历输入文本的每个字符

	for _, r := range text {
		// 条件判断：跳过空白、标点与符号

		if unicode.IsSpace(r) || unicode.IsPunct(r) || unicode.IsSymbol(r) {
			// 进入下一轮循环

			continue
			// 当前代码块结束

		}
		// 写入转小写后的字符

		b.WriteRune(unicode.ToLower(r))
		// 当前代码块结束

	}
	// 返回归一化结果给调用方，结束当前函数

	return b.String()
	// 当前代码块结束

}

// bigrams 生成中文按 rune 的二元组集合，用作相似度计算特征
// 定义函数 bigrams：返回二元组分词集合

func bigrams(text string) map[string]bool {
	// 将字符串转换为 rune 切片以正确处理中文

	runes := []rune(text)
	// 定义并初始化二元组集合

	set := make(map[string]bool)
	// 遍历相邻两个字符

	for i := 0; i+1 < len(runes); i++ {
		// 写入当前二元组

		set[string(runes[i:i+2])] = true
		// 当前代码块结束

	}
	// 返回集合给调用方，结束当前函数

	return set
	// 当前代码块结束

}

// similarity 使用二元组重叠系数计算两段文本相似度，取值 0~1
// 定义函数 similarity：返回归一化后的相似度

func similarity(a, b string) float64 {
	// 条件判断：任一输入为空时直接返回 0

	if a == "" || b == "" {
		// 返回 0 给调用方，结束当前函数

		return 0
		// 当前代码块结束

	}
	// 条件判断：一方包含另一方时视为完全匹配

	if strings.Contains(a, b) || strings.Contains(b, a) {
		// 返回 1 给调用方，结束当前函数

		return 1
		// 当前代码块结束

	}
	// 分别计算双方二元组集合

	setA := bigrams(a)
	// 计算乙方二元组集合

	setB := bigrams(b)
	// 条件判断：集合为空时执行下面逻辑

	if len(setA) == 0 || len(setB) == 0 {
		// 返回 0 给调用方，结束当前函数

		return 0
		// 当前代码块结束

	}
	// 定义并初始化命中计数

	hits := 0
	// 遍历甲方二元组集合

	for gram := range setA {
		// 条件判断：乙方包含当前二元组时命中数加一

		if setB[gram] {
			// 命中数加一

			hits++
			// 当前代码块结束

		}
		// 当前代码块结束

	}
	// 计算较小集合大小作为分母，返回重叠系数

	denom := len(setA)
	// 条件判断：乙方集合更小时更新分母

	if len(setB) < denom {
		// 为 denom 赋值 len(setB)

		denom = len(setB)
		// 当前代码块结束

	}
	// 返回重叠系数给调用方，结束当前函数

	return float64(hits) / float64(denom)
	// 当前代码块结束

}

// containsAny 判断归一化文本是否命中任一关键词
// 定义函数 containsAny：用于品牌类问题的关键词识别

func containsAny(text string, words []string) bool {
	// 遍历关键词列表

	for _, word := range words {
		// 条件判断：文本包含当前关键词时执行下面逻辑

		if strings.Contains(text, normalize(word)) {
			// 返回 true 给调用方，结束当前函数

			return true
			// 当前代码块结束

		}
		// 当前代码块结束

	}
	// 返回 false 给调用方，结束当前函数

	return false
	// 当前代码块结束

}

// emotionActionFor 依据真实 FAQ 条目 id 返回前端支持的情绪与动作取值
// 定义函数 emotionActionFor：返回情绪名与动作名

func emotionActionFor(id string) (string, string) {
	// 取条目 id 的末尾序号，兼容原有 q 系列与新角色 s/r/v 系列
	tail := id
	if len(id) > 0 {
		tail = id[len(id)-1:]
	}
	// 按序号区分情绪与动作，全部取自前端 emotionAction 与 Live2DAction 的真实枚举
	switch tail {
	case "1", "3":
		// 寄件/取件/咨询类使用积极情绪与点头动作
		return "positive", "nod"
	case "2":
		// 收费/信息类使用中性情绪与眨眼动作
		return "neutral", "blink"
	case "4":
		// 故障/报修类使用难过情绪与摇头动作
		return "sad", "shake"
	case "5":
		// 凭证/建议类使用开心情绪与微笑动作
		return "happy", "smile"
	}
	// 默认返回开心情绪与微笑动作
	return "happy", "smile"
}

// BuildReply 依据用户文本与角色构造严格对齐前端协议的 assistant 回复（角色间内容互不相同）@ticket 07
// 定义函数 BuildReply：返回可直接序列化的 AgentReply
// 参数 role 为账号判定的角色 sender/receiver/service，决定使用哪一套人设与专属 FAQ
func BuildReply(userContent string, role string) AgentReply {
	// 获取只读通用知识库（品牌介绍使用）
	k := getKnowledge()
	// 按角色获取专属知识，未识别角色回退到智能客服专员
	rk := getRoleKnowledge()[role]
	if rk == nil {
		rk = getRoleKnowledge()["service"]
	}
	// 取当前角色的专属 FAQ 列表
	faqs := rk.FAQs
	// 归一化用户输入
	q := normalize(userContent)
	// 定义并初始化回复文本
	text := ""
	// 定义并初始化情绪名
	emotion := "neutral"
	// 定义并初始化动作名
	actionName := "blink"
	// 条件判断：处理品牌/介绍类问题（回答中带入当前角色人设，体现角色差异）
	if containsAny(q, []string{"品牌", "介绍", "你们是谁", "八维", "简介", "关于", "你是谁", "你的身份"}) {
		// 品牌介绍前拼接当前角色人设，使三角色自我介绍互不相同
		persona := rk.Profile.Persona
		if persona == "" {
			persona = k.BrandIntro
		}
		text = persona + " " + k.BrandIntro + " " + strings.Join(k.BrandLines, " ")
		// 情绪设置为开心
		emotion = "happy"
		// 动作设置为微笑
		actionName = "smile"
	} else {
		// 定义并初始化最高分
		bestScore := -1.0
		// 定义并初始化最佳条目下标
		bestIndex := -1
		// 遍历当前角色专属 FAQ 条目计算匹配分
		for i, item := range faqs {
			// 计算用户问题与标准问题的相似度，并乘以条目权重
			score := similarity(q, normalize(item.Q)) * item.Weight
			// 遍历条目标签做关键词加权
			for _, tag := range item.Tags {
				// 条件判断：用户问题命中标签时加权
				if strings.Contains(q, normalize(tag)) {
					// 分数增加 0.3
					score += 0.3
				}
			}
			// 条件判断：刷新最高分
			if score > bestScore {
				// 更新最高分
				bestScore = score
				// 更新最佳条目下标
				bestIndex = i
			}
		}
		// 条件判断：命中质量达到阈值时使用角色专属 FAQ 答案
		if bestIndex >= 0 && bestScore >= 0.18 {
			// 回复文本使用角色专属真实 FAQ 答案，禁止臆造
			text = faqs[bestIndex].A
			// 依据条目 id 设置前端真实支持的情绪与动作
			emotion, actionName = emotionActionFor(faqs[bestIndex].ID)
		}
	}
	// 条件判断：未命中任何真实知识时给出兜底，兜底文案由角色人设与该角色真实问题列表组成
	if text == "" {
		// 收集当前角色 FAQ 的标准问题作为可问引导
		questions := make([]string, 0, len(faqs))
		// 遍历角色专属 FAQ 条目
		for _, item := range faqs {
			// 追加标准问题
			questions = append(questions, item.Q)
		}
		// 兜底文本由当前角色人设与真实问题列表组成，三角色互不相同
		persona := rk.Profile.Persona
		if persona == "" {
			persona = k.BrandIntro
		}
		text = persona + " 您可以问我：" + strings.Join(questions, "、")
		// 情绪设置为中性
		emotion = "neutral"
		// 动作设置为眨眼
		actionName = "blink"
	}
	// 构造并返回严格对齐前端协议的回复对象
	return AgentReply{
		// 设置角色为 assistant
		Role: "assistant",
		// 设置回复载荷
		Payload: AgentReplyPayload{
			// 设置真实回复文本
			Text: text,
			// 设置前端真实支持的情绪
			Emotion: emotion,
			// 使用场景默认值 chat
			Scene: "chat",
			// 设置表情动作联动
			Action: &Live2DAction{
				// 设置前端真实支持的动作名
				Action: actionName,
				// 设置对应情绪
				Emotion: emotion,
				// 设置前端 emotionAction 默认时长
				DurationMS: 1200,
			},
		},
	}
}
