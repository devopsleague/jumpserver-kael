<script setup>
</script>

<template>
  <div class="loading">
    <div class="box loading-box">
      <span></span>
      <span></span>
      <span></span>
      <span></span>
      <span></span>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.loading {
  position: fixed;
  width: 100%;
  height: 100%;
}
.loading-box{
  position: absolute;
  bottom: 134px;
  left: 40%;
  text-align: center;
  height: 15px;
  margin: 0 auto;
}
.loading-box span{
  display: inline-block;
  width: 15px;
  height: 100%;
  margin-right: 5px;
  border-radius: 50%;
  background: rgb(182, 189, 198);
  animation: load 1.04s ease infinite;
}
.loading-box span:last-child{
  margin-right: 0px;
}
@keyframes load{
  0%{
    opacity: 1;
  }
  100%{
    opacity: 0;
  }
}
.loading-box span:nth-child(1){
  animation-delay: 0.13s;
}
.loading-box span:nth-child(2){
  animation-delay: 0.26s;
}
.loading-box span:nth-child(3){
  animation-delay: 0.39s;
}
.loading-box span:nth-child(4){
 animation-delay: 0.52s;
}
.loading-box span:nth-child(5){
  animation-delay: 0.65s;
}
</style>