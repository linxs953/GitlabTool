# Gitalb Automation Cli Tool

## What is this

> ### gitlab协作的自动化工具，将一些常用的操作封装成命令行，减少时间，提升效率

## Why do this

> ### 将常用的gitlab操作封装成cli操作，比如查看cicd，快速创建issue，提高效率


```shell
## 通过模板新建一个issue
1. mt new [issue_title] [issue-template-filename] [assignee-name]

## 根据当前`project`打开对应的cicd地址
2. mt cicd 

## 更改这个配置文件已存在的key的值  
3. mt config [key] [value]

## 列出当前有权限的group和project
4. mt list [groups] / [projects] 

## 更改当前的group和project
5. mt sw group [group_name] / project [project_name] 

```