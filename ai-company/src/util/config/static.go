package config

const (
	MASTER_PROMPT_FILE_PATH = "/prompt/tasks"
	WORKER_PROMPT_FILE_PATH = "/prompt/worker_tasks"
	INSTRUCTION_FILE_NAME   = "instruction.md"

	ENV_WORKER_NUM = "WORKER_NUM"
	ENV_WORKER_ID  = "WORKER_ID"
	ENV_QUEUE_HOST = "QUEUE_HOST"

	WORKER_FILE_NAME = "worker-%d.md"
)
