package golang

var PipelineBuilder = `name: Pipeline Builder
on:
  push:
    branches: [ master, main ]
  pull_request:
    branches: [ master, main ]
jobs:
  [[- if .Build.Enable]]
  build:
    runs-on: ubuntu-latest
    steps:
    - name: Build
      run: [[.Build.Command]]
  [[- end]]
  [[- if .Test.Enable]]
  test:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v2
    - name: Set up Go
      uses: actions/setup-go@v2
      with:
        go-version: 1.17
    - name: Test
      run: [[.Test.Command]] [[- if .Test.Coverage.Enable]] -race -covermode=atomic -coverprofile=[[.Test.Coverage.Output]] [[- end]]
#    TODO(ironcore): this is a bug. If it runs on main branch, it should not comment PR since there isn't any.
#    [[- if .Test.Coverage.Enable]]
#    - name: comment PR
#      uses: machine-learning-apps/pr-comment@master
#      env:
#        GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
#      with:
#        path: [[.Test.Coverage.Output]]
#    [[- end]]
  [[- end]]
  tag:
    name: Tag
    needs: [test]
    if: ${{ github.event_name == 'push' }}
    runs-on: ubuntu-latest
    outputs:
      new_tag: ${{ steps.tag_version.outputs.new_tag }}
    steps:
      - uses: actions/checkout@v2
      - name: Bump version and push tag
        id: tag_version
        uses: mathieudutour/github-tag-action@v5.6
        with:
          github_token: ${{ secrets.GITHUB_TOKEN }}
          tag_prefix: ""
  [[- if .Docker.Enable]]
  image:
    name: Build Docker Image
    needs: [tag]
    if: ${{ github.event_name == 'push' }}
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - name: Login to Docker Hub
        uses: docker/login-action@v1
        with:
          username: ${{ secrets.DOCKERHUB_USERNAME }}
          password: ${{ secrets.DOCKERHUB_TOKEN }}
      - name: Build and push
        id: docker_build
        uses: docker/build-push-action@v2
        with:
          push: true
          tags: ${{ secrets.DOCKERHUB_USERNAME }}/[[.Docker.Repo]]:${{needs.tag.outputs.new_tag}}
  [[- end]]
`
