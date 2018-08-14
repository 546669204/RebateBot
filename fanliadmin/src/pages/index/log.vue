<template>
    <div>
        <div v-for="(val,index) in '1234567891'" :key="index" class="icon-list">
            
                <i class="el-icon-service"></i> 
            
            
            <div>
                <span>Semantic-Org/Semantic-UI</span>
                <span>Updated 10 mins ago</span>
            </div>
            
        </div>
    </div>
</template>
<style>
.icon-list{
    display: flex;
    flex-direction: row;
    align-content: center;

}
.icon-list *{
    display: flex;
}
.icon-list>div{
    flex-direction: column;
}
.icon-list>i{
    align-items: center;
    font-size: 32px;
}
</style>

<script>
export default {
    name:"IndexLog",
    data(){
        return {
            
        }
    },
    methods:{

    }
}    
</script>